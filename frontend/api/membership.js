import httpRequest from "~/utils/request";

export const membershipApi = {
  /** 获取套餐列表（公开） */
  async getPlans(params = {}, server = true) {
    let options = { url: `/api/v1/membership/plans`, method: "get", params, server };
    return await httpRequest.post(options);
  },

  /** 升级会员 */
  async upgrade(params = {}, server = true) {
    let options = { url: `/api/v1/membership/upgrade`, method: "post", params, server };
    return await httpRequest.post(options);
  },

  /** 购买积分 */
  async buyCredits(params = {}, server = true) {
    let options = { url: `/api/v1/membership/buy-credits`, method: "post", params, server };
    return await httpRequest.post(options);
  },
};
