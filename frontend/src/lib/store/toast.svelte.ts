type ToastType = {
  message: string;
  type?: 'success' | 'error' | 'info';
  duration?: number; // Duration in milliseconds
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left';
  id: number;
};

class ToastStore {
  private _toasts: ToastType[] = $state([]);
  private _defaultDuration = 3000;
  private _nextId = 0;

  add(
    msg: string,
    type: ToastType['type'] = 'success',
    duration: ToastType['duration'] = this._defaultDuration,
    position: ToastType['position'] = 'top-right'
  ) {
    this._nextId += 1;
    const id = this._nextId;
    this._toasts.push({
      message: `${id}: ${msg}`,
      id,
      type,
      duration,
      position,
    });

    const timer = () => {
      this.remove(id);
    };

    window.setTimeout(timer, duration);
  }

  setDuration(duration: number) {
    this._defaultDuration = duration;
  }

  get all() {
    return this._toasts;
  }

  remove(id: number) {
    this._toasts = this._toasts.filter(t => t.id !== id);
  }

  pop() {
    if (this._toasts.length < 1) {
      this._toasts = [];
      return;
    }

    this._toasts = this._toasts.slice(1);
  }

  get duration() {
    return this._defaultDuration;
  }
}

export const toastStore = new ToastStore();
