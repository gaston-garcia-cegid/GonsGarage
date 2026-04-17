/**
 * Runs before React paint (Next `beforeInteractive`) to avoid theme flash.
 * Keep in sync with `theme.ts` (storage key + resolution rules).
 */
export const THEME_INLINE_SCRIPT = `(function(){try{var k='gonsgarage-theme';var raw=localStorage.getItem(k)||'system';var r='light';if(raw==='dark'||raw==='light'){r=raw;}else{r=window.matchMedia('(prefers-color-scheme:dark)').matches?'dark':'light';}document.documentElement.setAttribute('data-theme',r);document.documentElement.setAttribute('data-theme-pref',raw);}catch(e){}})();`;
