// Prevent images from being dragged or selected programmatically
if (typeof document !== 'undefined') {
  function disableImageDragAndSelect() {
    try {
      const imgs = document.querySelectorAll('img');
      imgs.forEach(img => {
        try { img.setAttribute('draggable', 'false'); } catch (e) {}
      });
      document.addEventListener('dragstart', (e) => {
        if (e.target && e.target.tagName === 'IMG') e.preventDefault();
      }, { capture: true });
      // Prevent right-click context menu on prominent images (hero / feature images)
      document.addEventListener('contextmenu', (e) => {
        try {
          const t = e.target;
          if (!t) return;
          if (t.tagName === 'IMG') {
            // keep it scoped to hero/feature/logo images to avoid breaking site-wide UX
            if (t.closest('.hero') || t.classList.contains('featureSvg_GfXr') || t.classList.contains('themedComponent_mlkZ')) {
              e.preventDefault();
            }
          }
        } catch (err) {}
      }, { capture: true });
    } catch (e) {
      // ignore
    }
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', disableImageDragAndSelect);
  } else {
    disableImageDragAndSelect();
  }
}

export default {};
