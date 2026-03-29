
import {baseUrl} from "~/config/api"
import httpRequest from "~/utils/request"
export const identifyApi = {
  sendFile: async (params: any) => {
    console.log(params);
    let url = baseUrl.identifySend
    let reslut = await httpRequest.post(url,params)
    // const response = await fetch('/api/ajaxData', {
    //     method: 'POST',
    //     headers: { 'Content-Type': 'multipart/form-data' },
    //     body: formData
    // })
    // const data = await response.json()
    // return data
  },
};
