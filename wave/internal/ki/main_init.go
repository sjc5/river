package ki

import (
	"encoding/json"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/river-now/river/kit/colorlog"
	"github.com/river-now/river/kit/safecache"
	"golang.org/x/sync/semaphore"
)

const __internal_full_dev_reset_less_go_mrkr = "__internal_full_dev_reset_less_go_mrkr"

/////////////////////////////////////////////////////////////////////
/////// MAIN INIT
/////////////////////////////////////////////////////////////////////

type MainInitOptions struct {
	IsDev     bool
	IsRebuild bool
}

func (c *Config) MainInit(opts MainInitOptions, calledFrom string) {
	// LOGGER
	if c.Logger == nil {
		c.Logger = colorlog.New("wave")
	}

	c.fileSemaphore = semaphore.NewWeighted(100)

	// USER CONFIG
	c._uc = new(UserConfig)
	if err := json.Unmarshal(c.ConfigBytes, c._uc); err != nil {
		c.panic("failed to unmarshal user config", err)
	}

	// CLEAN SOURCES
	c.cleanSources = CleanSources{
		Dist:          filepath.Clean(c._uc.Core.DistDir),
		PrivateStatic: filepath.Clean(c._uc.Core.StaticAssetDirs.Private),
		PublicStatic:  filepath.Clean(c._uc.Core.StaticAssetDirs.Public),
	}
	if c._uc.Core.CSSEntryFiles.Critical != "" {
		c.cleanSources.CriticalCSSEntry = filepath.Clean(c._uc.Core.CSSEntryFiles.Critical)
	}
	if c._uc.Core.CSSEntryFiles.NonCritical != "" {
		c.cleanSources.NonCriticalCSSEntry = filepath.Clean(c._uc.Core.CSSEntryFiles.NonCritical)
	}

	// DIST LAYOUT
	c._dist = toDistLayout(c.cleanSources.Dist)

	c.InitRuntimeCache()

	// AFTER HERE, ALL DEV-TIME STUFF
	if !opts.IsDev {
		return
	}

	c.dev.mu.Lock()
	defer c.dev.mu.Unlock()

	c.kill_browser_refresh_mux()

	c._rebuild_cleanup_chan = make(chan struct{})

	c.cleanWatchRoot = filepath.Clean(c._uc.Watch.WatchRoot)

	SetModeToDev()

	// HEALTH CHECK ENDPOINT
	if c._uc.Watch.HealthcheckEndpoint == "" {
		c._uc.Watch.HealthcheckEndpoint = "/"
	}

	if !opts.IsRebuild {
		c.browserTabManager = newClientManager()
		go c.browserTabManager.start()
	}

	c.ignoredFilePatterns = []string{
		c.get_binary_output_path(),
	}

	c.naiveIgnoreDirPatterns = []string{
		"**/.git",
		"**/node_modules",
		c._dist.S().Static.FullPath(),
		filepath.Join(c.cleanSources.PublicStatic, noHashPublicDirsByVersion[0]),
		filepath.Join(c.cleanSources.PublicStatic, noHashPublicDirsByVersion[1]),
	}

	for _, p := range c.naiveIgnoreDirPatterns {
		c.ignoredDirPatterns = append(c.ignoredDirPatterns, filepath.Join(c.cleanWatchRoot, p))
	}
	for _, p := range c._uc.Watch.Exclude.Dirs {
		c.ignoredDirPatterns = append(c.ignoredDirPatterns, filepath.Join(c.cleanWatchRoot, p))
	}
	for _, p := range c._uc.Watch.Exclude.Files {
		c.ignoredFilePatterns = append(c.ignoredFilePatterns, filepath.Join(c.cleanWatchRoot, p))
	}

	c.defaultWatchedFiles = []WatchedFile{
		{
			Pattern:       filepath.Join(c.cleanSources.PublicStatic, "**/*"),
			OnChangeHooks: []OnChangeHook{{Cmd: __internal_full_dev_reset_less_go_mrkr}},
		},
	}

	includeDefaults := c._uc.River != nil
	if c._uc.River != nil && c._uc.River.IncludeDefaults != nil && !*c._uc.River.IncludeDefaults {
		includeDefaults = false
	}

	if includeDefaults {
		relClientRouteDefsFile, err := filepath.Rel(c.cleanWatchRoot, c._uc.River.ClientRouteDefsFile)
		if err != nil {
			c.panic("failed to get relative path for ClientRouteDefsFile", err)
		}

		c.defaultWatchedFiles = append(c.defaultWatchedFiles, WatchedFile{
			Pattern:       filepath.Join(c.cleanSources.PrivateStatic, c._uc.River.HTMLTemplateLocation),
			OnChangeHooks: []OnChangeHook{{Cmd: __internal_full_dev_reset_less_go_mrkr}},
		})

		c.defaultWatchedFiles = append(c.defaultWatchedFiles, WatchedFile{
			Pattern:       filepath.ToSlash(relClientRouteDefsFile),
			OnChangeHooks: []OnChangeHook{{Cmd: __internal_full_dev_reset_less_go_mrkr}},
		})

		c.defaultWatchedFiles = append(c.defaultWatchedFiles, WatchedFile{
			Pattern:       "**/*.go",
			OnChangeHooks: []OnChangeHook{{Cmd: "DevBuildHook", Timing: "concurrent"}},
		})

		relHTMLTemplateLocation, err := filepath.Rel(c.cleanWatchRoot, c._uc.River.HTMLTemplateLocation)
		if err != nil {
			c.panic("failed to get relative path for HTMLTemplateLocation", err)
		}

		c.defaultWatchedFiles = append(c.defaultWatchedFiles, WatchedFile{
			Pattern:    filepath.ToSlash(relHTMLTemplateLocation),
			RestartApp: true,
		})

		relTSGenOutPath, err := filepath.Rel(c.cleanWatchRoot, c._uc.River.TSGenOutPath)
		if err != nil {
			c.panic("failed to get relative path for TSGenOutPath", err)
		}

		c.ignoredFilePatterns = append(
			c.ignoredFilePatterns,
			filepath.ToSlash(relTSGenOutPath),
		)
	}

	// Loop through all WatchedFiles...
	for i, wfc := range c._uc.Watch.Include {
		// and make each WatchedFile's Pattern relative to cleanWatchRoot...
		c._uc.Watch.Include[i].Pattern = filepath.Join(c.cleanWatchRoot, wfc.Pattern)
		// then loop through such WatchedFile's OnChangeHooks...
		for j, oc := range wfc.OnChangeHooks {
			// and make each such OnChangeCallback's ExcludedPatterns also relative to cleanWatchRoot
			for k, p := range oc.Exclude {
				c._uc.Watch.Include[i].OnChangeHooks[j].Exclude[k] = filepath.Join(c.cleanWatchRoot, p)
			}
		}
	}

	c.matchResults = safecache.NewMap(c.get_initial_match_results, c.match_results_key_maker, nil)

	if c.watcher != nil {
		if err := c.watcher.Close(); err != nil {
			c.panic("failed to close watcher", err)
		}
		c.watcher = nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		c.panic("failed to create watcher", err)
	}

	c.watcher = watcher

	if err := c.add_directory_to_watcher(c.cleanWatchRoot); err != nil {
		c.panic("failed to add directory to watcher", err)
	}
}
