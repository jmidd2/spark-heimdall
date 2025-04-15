<!-- frontend/src/lib/components/Toast.svelte -->
<script lang="ts">
  type ToastProps = {
    message: string;
    type?: 'success' | 'error' | 'info';
    duration?: number; // Duration in milliseconds
    position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left';
  }

  // Props
  const { message, type = 'success', duration = 3000, position = 'top-right' }: ToastProps = $props();

  // State
  let isVisible = $state(true);

  // Automatically hide the toast after the specified duration
  $effect(() => {
    const timer = setTimeout(() => {
      isVisible = false;
    }, duration);

    return () => clearTimeout(timer);
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
        <button class="toast-icon">
            {@html getIcon()}
        </button>
        <div class="toast-content">
            {message}
        </div>
        <button class="toast-close" onclick={close}>
            <span aria-hidden="true">Close</span>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
        </button>
    </div>
{/if}

<style>
    .toast {
        position: fixed;
        display: flex;
        align-items: center;
        padding: 0.75rem 1rem;
        border-radius: 4px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        animation: slideIn 0.3s ease-out;
        z-index: 1000;
        min-width: 250px;
        max-width: 450px;
    }

    /* Position variants */
    .toast-top-right {
        top: 1rem;
        right: 1rem;
    }

    .toast-top-left {
        top: 1rem;
        left: 1rem;
    }

    .toast-bottom-right {
        bottom: 1rem;
        right: 1rem;
    }

    .toast-bottom-left {
        bottom: 1rem;
        left: 1rem;
    }

    /* Type variants */
    .toast-success {
        background-color: #1a3e29;
        border-left: 4px solid #27ae60;
        color: #b9ffc9;
    }

    .toast-error {
        background-color: #482121;
        border-left: 4px solid #e74c3c;
        color: #f5c2c2;
    }

    .toast-info {
        background-color: #1e3a5a;
        border-left: 4px solid #3498db;
        color: #a8d7ff;
    }

    .toast-icon {
        margin-right: 0.75rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .toast-content {
        flex: 1;
    }

    .toast-close {
        background: none;
        border: none;
        color: inherit;
        opacity: 0.7;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 0;
        margin-left: 0.75rem;
    }

    .toast-close:hover {
        opacity: 1;
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
            width: calc(100% - 2rem);
            max-width: none;
        }

        .toast-top-right,
        .toast-top-left {
            top: 0.5rem;
            right: 0.5rem;
            left: 0.5rem;
        }

        .toast-bottom-right,
        .toast-bottom-left {
            bottom: 0.5rem;
            right: 0.5rem;
            left: 0.5rem;
        }
    }
</style>