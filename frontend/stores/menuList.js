import { defineStore } from "pinia";

export const useMenuListStore = defineStore("menuList", {
  state: () => ({
    navItems: [
      { icon: "fa-solid fa-house", label: "首页", to: "/" },
      { icon: "fa-solid fa-code", label: "代码生成", to: "/code" },
      { icon: "fa-regular fa-folder", label: "任务列表", to: "/tasks" },
      { icon: "fa-regular fa-copy", label: "模板市场", to: "/templates" },
    ],
  }),

  actions: {},

  getters: {
    /**
     * 含 active 字段的导航列表，当前路径匹配时 active=true。
     * 当前路径必须由调用方传入：getter 可能在渲染/副作用等非 setup 上下文执行，
     * 在这里调 useRouter() 拿不到实例（浏览器报 reading 'currentRoute' of undefined）。
     */
    navItemsWithActive: (state) => (currentPath) =>
      state.navItems.map((item) => ({
        ...item,
        active: currentPath === item.to,
      })),
  },
});
