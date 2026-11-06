import httpRequest from "~/utils/request";

export const uiProjectApi = {
  // 整页 AI 描述：主图 URL + 整页区域树 → { tree }（与请求同构，每个节点 summary 已填 AI 描述）
  async describeRegion(params = {}, server = true) {
    let options = {
      url: `/api/v1/ui-project/describe`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
};
