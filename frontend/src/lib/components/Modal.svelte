<script lang="ts">
import type { Snippet } from 'svelte';

type ModalProps = {
  onClose: () => void;
  showModal: boolean;
  children: Snippet;
  onSubmit: () => Promise<void>;
  title: string;
};

// Props
let {
  onClose,
  showModal = $bindable(),
  children,
  onSubmit,
  title,
}: ModalProps = $props();

let dialog = $state<HTMLDialogElement | undefined>();

$effect(() => {
  if (showModal) dialog?.showModal();
  if (!showModal) dialog?.close();
});

function handleClose() {
  showModal = false;
  dialog?.close();
  onClose();
}

let isSaving = $state(false);

async function handleSubmit(e: SubmitEvent) {
  e.preventDefault();
  isSaving = true;
  await onSubmit();
  isSaving = false;
}
</script>

<dialog bind:this={dialog} class="modal"
        onclose={() => (showModal = false)}
        onclick={(e) => { if (e.target === dialog) dialog.close(); }}>
    <div>
        <div class="modal-header">
            <h2>{title}</h2>
            <button type="button" class="close-btn" onclick={handleClose}>&times;</button>
        </div>
        <div class="modal-content">
            <form onsubmit={handleSubmit}>
                {@render children()}

                <div class="form-actions">
                    <div class="btn-group">
                        <button
                                type="submit"
                                class="btn btn-primary"
                                disabled={isSaving}
                        >
                            {isSaving ? 'Saving...' : 'Save'}
                        </button>
                        <button type="button" class="btn btn-secondary" onclick={handleClose}>
                            Cancel
                        </button>
                    </div>
                </div>
            </form>
        </div>
    </div>
</dialog>

<style>
    dialog {
        border-radius: var(--heimdall-rounded);
        border: none;
        padding: 0;
        margin: auto;
        background-color: var(--heimdall-bg-modal);
        min-width: var(--heimdall-modal-width);
        max-width: var(--heimdall-modal-width);
        height: var(--heimdall-modal-height);

        &::backdrop {
            background-color: rgba(0, 0, 0, 0.5);
        }

        .modal-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: var(--heimdall-spacing-lg);
            border-bottom: 1px solid var(--heimdall-border-color);
        }

        .modal-header h2 {
            margin: 0;
            color: var(--heimdall-text-heading);
        }

        .close-btn {
            background: none;
            border: none;
            font-size: var(--heimdall-font-size-2xl);
            color: var(--heimdall-text-secondary);
            cursor: pointer;
        }

        form {
            padding: var(--heimdall-spacing-lg);
        }

        .form-actions {
            border-top: 1px solid var(--heimdall-border-color);
            display: block;
            gap: var(--heimdall-spacing-sm);

            & > :not(:last-child) {
                border-top-width: 0;
                border-bottom-width: 1px;
                border-color: var(--heimdall-border-color);
                border-style: solid;
                border-left-style: none;
                border-right-style: none;
            }

            .btn-group {
                padding-top: var(--heimdall-spacing-md);
                padding-bottom: var(--heimdall-spacing-md);
            }
        }
    }
</style>