
import httpRequest from "~/utils/request"

export const identifyApi = {
 async sendFile(params = {}, server = true) {
    let options = {
        url: `/api/generate/send`,
        method: "post",
        params,
        server,
    }
    let result = await httpRequest.post(options);
    return result;
  },
}