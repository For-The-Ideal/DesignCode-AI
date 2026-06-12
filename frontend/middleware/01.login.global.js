import { useUserStore } from "~/stores/user";
export default defineNuxtRouteMiddleware(async (to, from) => {
  try {
    const userStore = useUserStore();
    userStore.initialize();
  } catch (error) {
    console.error("Error initializing store:", error);
  }
  return true;
});
