/**
 * Composables index
 * Centralized export of all composables
 */

// Phase 1-3: Infrastructure, Performance, Responsive
export { useBreakpoint, useMediaQuery, useResponsiveValue } from './useBreakpoint'
export { useVirtualScroll, useDynamicVirtualScroll } from './useVirtualScroll'
export { useTouchGestures, useSwipeGestures, usePinchGesture, useLongPress, useTouchDrag } from './useTouchGestures'

// Phase 4: Component Refactoring
export { useGanttKeyboard, useShortcutPanel } from './useGanttKeyboard'
export { useGanttTooltip, useGanttTooltipAdvanced } from './useGanttTooltip'
export { useGanttSelection, useKeyboardNavigation } from './useGanttSelection'

// Phase 5: UX Optimization
export { useTheme, THEME_PRESETS, getThemeIcon, getThemeElementPlusIcon } from './useTheme'
export {
  useAria,
  useFocusManager,
  useScreenReader,
  useKeyboardNav,
  useHighContrast,
  useReducedMotion
} from './useA11y'
