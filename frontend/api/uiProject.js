import httpRequest from "~/utils/request";

/**
 * UI 图转代码 —— 项目相关接口
 * 走 SSR 代理（utils/request.js），与 task/common 同一套返回约定（code === 200 为成功）。
 */
export const uiProjectApi = {

  // 顶层区域识别：进编辑器时按图片调 AI（后端用 Qwen-VL / capability=read）
  // 返回规范化的顶层区域目录，供框选顶层时按 name 命中自动填充描述。
  // 返回体约定：{ code: 200, data: { regions: [{ name, summary, styleHints, behavior, childrenHint, layoutHint }] } }
  async analyzeRegions(imageUrl, extraOptions = {},server = true) {
    let options = {
      url: `/api/v1/ui-project/analyze`,
      method: "post",
      params: {
        imageUrl,
      },
      server,
      ...extraOptions,
    };
    let result = await httpRequest.post(options);
    return result;
  },

};
