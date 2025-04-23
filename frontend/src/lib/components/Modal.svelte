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