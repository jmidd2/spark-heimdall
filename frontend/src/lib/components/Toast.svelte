<!-- frontend/src/lib/components/Toast.svelte -->
<script lang="ts">
type ToastProps = {
  message: string;
  type?: 'success' | 'error' | 'info';
  duration?: number; // Duration in milliseconds
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left';
};

// Props
const {
  message,
  type = 'success',
  duration = 3000,
  position = 'top-right',
}: ToastProps = $props();

// State
let isVisible = $state(true);

// Automatically hide the toast after the specified duration
$effect(() => {
  if (duration > 0) {
    const timer = setTimeout(() => {
      isVisible = false;
    }, duration);

    return () => clearTimeout(timer);
  }
});

// Close the toast
function close() {
  isVisible = false;
}

// Get icon based on toast type
function getIcon() {
  switch (type) {
    case 'success':
      return `
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
            <polyline points="22 4 12 14.01 9 11.01"></polyline>
          </svg>
        `;
    case 'error':
      return `
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="15" y1="9" x2="9" y2="15"></line>
            <line x1="9" y1="9" x2="15" y2="15"></line>
          </svg>
        `;
    default:
      return `
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="16" x2="12" y2="12"></line>
            <line x1="12" y1="8" x2="12.01" y2="8"></line>
          </svg>
        `;
  }
}
</script>

{#if isVisible}
    <div class="toast toast-{type} toast-{position}">
        <div class="toast-icon">
            {@html getIcon()}
        </div>
        <div class="toast-content">
            {message}
        </div>
        <button class="toast-close" onclick={close}>
            <span class="hidden" aria-hidden="true">Close</span>
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                 stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
        </button>
    </div>
{/if}

<style>
    .toast {
        line-height: 1rem;
        position: fixed;
        padding: var(--heimdall-spacing-md) var(--heimdall-spacing-lg) calc(var(--heimdall-spacing-md) - 0.1rem);
        border-radius: var(--heimdall-rounded);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        animation: slideIn 0.3s ease-out;
        z-index: 1000;
        min-width: 250px;
        max-width: 500px;
        display: grid;
        grid-template-columns: 20px 1fr auto;
        justify-content: space-around;
        align-content: center;
        align-items: center;
        vertical-align: middle;
    }

    .toast-icon {
        padding: 0;
        margin: 0.1rem 0 0 0;
    }

    .toast-content {
        padding: 0;
        margin: 0 var(--heimdall-spacing-md);
        text-align: center;
    }

    .toast-close {
        background: none;
        border: none;
        color: inherit;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        opacity: 0.7;
        cursor: pointer;
        padding: 0 calc(var(--heimdall-spacing-xs) - 0.05rem);
        margin: 0;
        border-radius: 10px;
    }

    .toast-close:hover {
        opacity: 1;
        background-color: rgba(255,255,255,0.2);
    }

    /* Position variants */
    .toast-top-right {
        top: var(--heimdall-spacing-lg);
        right: var(--heimdall-spacing-lg);
    }

    .toast-top-left {
        top: var(--heimdall-spacing-lg);
        left: var(--heimdall-spacing-lg);
    }

    .toast-bottom-right {
        bottom: var(--heimdall-spacing-lg);
        right: var(--heimdall-spacing-lg);
    }

    .toast-bottom-left {
        bottom: var(--heimdall-spacing-lg);
        left: var(--heimdall-spacing-lg);
    }

    /* Type variants */
    .toast-success {
        background-color: var(--heimdall-toast-success-bg);
        border-left: 4px solid var(--heimdall-toast-success-border);
        color: var(--heimdall-toast-success-text);
    }

    .toast-error {
        background-color: var(--heimdall-toast-error-bg);
        border-left: 4px solid var(--heimdall-toast-error-border);
        color: var(--heimdall-error-text);
    }

    .toast-info {
        background-color: var(--heimdall-toast-info-bg);
        border-left: 4px solid var(--heimdall-toast-info-border);
        color: var(--heimdall-toast-info-text);
    }

    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }

    /* Responsive adjustments */
    @media (max-width: 576px) {
        .toast {
            width: calc(100% - var(--heimdall-spacing-2xl));
            max-width: none;
        }

        .toast-top-right,
        .toast-top-left {
            top: var(--heimdall-spacing-sm);
            right: var(--heimdall-spacing-sm);
            left: var(--heimdall-spacing-sm);
        }

        .toast-bottom-right,
        .toast-bottom-left {
            bottom: var(--heimdall-spacing-sm);
            right: var(--heimdall-spacing-sm);
            left: var(--heimdall-spacing-sm);
        }
    }
</style>