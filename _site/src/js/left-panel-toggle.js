(function() {
  if (typeof window === 'undefined' || typeof document === 'undefined') return;

  function safe(fn) {
    try { fn(); } catch (e) { console.warn('left-panel-toggle error', e); }
  }

  function init() {
    // Toggle panel on click of the Menu link
    document.querySelectorAll('.navbar__link.left-panel-dropdown').forEach(function (anchor) {
      anchor.addEventListener('click', function (e) {
        e.preventDefault();
        var parent = anchor.closest('.left-panel-dropdown');
        if (!parent) return;
        parent.classList.toggle('dropdown--show');
      });
    });

    // Close any open left-panel when clicking outside
    document.addEventListener('click', function (e) {
      document.querySelectorAll('.left-panel-dropdown.dropdown--show').forEach(function (open) {
        if (!open.contains(e.target)) open.classList.remove('dropdown--show');
      });
    });

    // Close on ESC
    document.addEventListener('keydown', function (e) {
      if (e.key === 'Escape') {
        document.querySelectorAll('.left-panel-dropdown.dropdown--show').forEach(function (open) { open.classList.remove('dropdown--show'); });
      }
    });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function () { safe(init); });
  } else {
    safe(init);
  }
})();
