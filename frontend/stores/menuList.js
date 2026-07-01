import { defineStore } from "pinia";

export const useMenuListStore = defineStore("menuList", {
  state: () => ({
    navItems: [
      { icon: "fa-solid fa-code", label: "代码生成", to: "/code" },
      { icon: "fa-regular fa-folder", label: "任务列表", to: "/tasks" },
      { icon: "fa-regular fa-copy", label: "模板市场", to: "/templates" },
      { icon: "fa-regular fa-file", label: "我的项目", to: "/projects" },
    ],
  }),

  actions: {},

  getters: {
    getNavItems: (state) => state.navItems,
  },
});
