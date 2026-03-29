import { defineStore } from "pinia";

interface ToastParams {
  title: string;
  message: string;
  type: string;
}

export const useToastNotificationStore = defineStore("toastNotification", {
  state: () => ({
    visible: false,
    duration: 3000,
    toastParams: {
      title: "",
      message: "",
      type: "",
    },
  }),

  actions: {
    setToast(value: boolean, params: ToastParams, duration: number = 3000) {
      this.visible = value;
      this.duration = duration;
      this.toastParams = params;
    },

    hideToast() {
      this.visible = false;
      this.toastParams = {
        title: "",
        message: "",
        type: "",
      };
    },

    success(message: string) {
      this.setToast(true, { title: "成功", message, type: "success" });
    },

    error(message: string) {
      this.setToast(true, { title: "错误", message, type: "error" });
    },

    warning(message: string) {
      this.setToast(true, { title: "温馨提示", message, type: "warning" });
    },

    loading(){
      this.setToast(true, { title: "加载中", message: "加载稍后...", type: "loading" });
    },
    info(message: string) {
      this.setToast(true, { title: "提示", message, type: "info" });
    },
  },
});
