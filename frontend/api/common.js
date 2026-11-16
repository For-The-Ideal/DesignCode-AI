import httpRequest from "~/utils/request";

export const commonApi = {
  // 获取模板
  async getTemplate(params = {}, server = true) {
    // 不带 id → 请求 /api/v1/template，由后端随机返回一条
    let url = '/api/v1/template' + (params.id ? `/${params.id}` : '');
    let options = {
      url: url,
      method: "get",
      params:{},
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },

  // 上传图片到 COS（base64）
  // extraOptions 可透传 signal —— 配合 utils/cancel.js 的取消作用域使用
  async uploadImage(base64Data, filename, extraOptions = {}) {
    let options = {
      url: `/api/v1/upload`,
      method: "post",
      params: {
        image: base64Data,
        filename: filename || 'image.png',
      },
      server: false,
      ...extraOptions,
    };
    let result = await httpRequest.post(options);
    return result;
  },

};
