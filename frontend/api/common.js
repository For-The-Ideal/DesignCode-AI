import httpRequest from "~/utils/request";

export const commonApi = {
  // 获取模板
  async getTemplate(params = {}, server = true) {
    let options = {
      url: `/api/v1/template/${params.id}`,
      method: "get",
      params:{},
      server,
    };
    let result = await httpRequest.get(options);
    return result;
  },

  // 查询单个任务状态与结果
  async getTaskById(taskId, server = true) {
    let options = {
      url: `/api/v1/task/${taskId}`,
      method: "get",
      server,
    };
    let result = await httpRequest.get(options);
    return result;
  },

  // 查询当前用户任务状态统计
  async getUserTasks(server = true) {
    let options = {
      url: `/api/v1/task/status`,
      method: "get",
      server,
    };
    let result = await httpRequest.get(options);
    return result;
  },

  // 上传图片到 COS（base64）
  async uploadImage(base64Data, filename) {
    let options = {
      url: `/api/v1/upload`,
      method: "post",
      params: {
        image: base64Data,
        filename: filename || 'image.png',
      },
      server: false,
    };
    let result = await httpRequest.post(options);
    return result;
  },

  async generateUi(params = {}, server = true) {
    let options = {
      url: `/api/v1/task/create`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
};
