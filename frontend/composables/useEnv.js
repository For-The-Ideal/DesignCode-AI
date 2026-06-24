export const useEnv = () => {
  const config = useRuntimeConfig();
  return {
    apiBase: config.public.apiBase,
    wsBase: config.public.wsBase,
    cryptoKey: config.public.cryptoKey,
    cryptoIv: config.public.cryptoIv,
    buildId: config.app.buildId,
  };
};
