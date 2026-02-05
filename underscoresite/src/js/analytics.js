// Analytics shim removed per user request.
// Provide no-op exports so other modules importing analytics do not fail.
export default function initAnalytics(){ /* noop */ }
export function trackEvent(){ /* noop */ }
export function isEnabled(){ return false }
