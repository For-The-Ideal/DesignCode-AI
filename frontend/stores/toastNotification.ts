import { defineStore } from "pinia";

interface ToastParams {
  title: string;
  message: string;
}

export const useToastNotificationStore = defineStore("toastNotification", {
  state: () => ({
    showToast: false,
    toastParams: {
      title: "",
      message: "",
    },
  }),

  actions: {
    setShowToast(value: boolean, params: ToastParams) {
      this.showToast = value;
      this.toastParams = params;
    },
  },
});
