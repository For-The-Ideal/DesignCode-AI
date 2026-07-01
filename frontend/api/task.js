import httpRequest from "~/utils/request";

export const taskApi = {

  // 查询单个任务
  async getTaskById(taskId, server = true) {
    let options = {
      url: `/api/v1/task/${taskId}`,
      method: "get",
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },

  // 查询当前用户任务状态统计
  async getTaskStatus(server = true) {
    let options = {
      url: `/api/v1/task/status`,
      method: "get",
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },

  // 创建任务
  async taskCreate(params = {}, server = true) {
    let options = {
      url: `/api/v1/task/create`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },

  // 查询任务列表
  async getTaskList(server = true) {
    let options = {
      url: `/api/v1/task/list`,
      method: "post",
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },

  // 删除任务
  async deleteTask(taskId, server = true) {
    let options = {
      url: `/api/v1/task/${taskId}`,
      method: "delete",
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
};
