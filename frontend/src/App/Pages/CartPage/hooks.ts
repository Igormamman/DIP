import saveAs from "file-saver"
import { useContext, useEffect, useState } from "react"
import { useApi } from "../../../api/api"
import { AppContext, store } from "../../AppContext"
import { Store } from "react-notifications-component"
import axios from "axios"

export const useCartPage = () => {

    const [items, setItems] = useState<string[]>()
    const { api } = useApi()
    const [loadState, setLoadState] = useState(0)
    const [loading, setLoading] = useState(false)
    const { userStore } = useContext(AppContext)

    useEffect(() => {
        getItems()
    }, [])

    const getItems = () => {
        let itemString = window.localStorage.getItem("SelectedPhotos")
        if (itemString != null) {
            setItems(JSON.parse(itemString))
        }
    }

    const [job, setJob] = useState([]);
    const [imagesUploaded, setImagesUploaded] = useState(null);

    //initialize jsZip
    var JSZip = require("jszip");

    const DownloadFileFromS3 = async (UID: string, fileName: string, zip: any) => {
        let getURL = `https://photo.marshalone.ru/api/photo/file/get?UUID=${UID}&type=`
        await api.get(getURL, { withCredentials: false }).then(
            async (response) => {
                if (response.status < 400) {
                    let presignedURl: string = response.data["URL"]
                    if (presignedURl != null) {
                        await api.get(presignedURl, { withCredentials: false, responseType: 'blob' }).then(
                            (innerResponse) => {
                                let blob = new Blob([innerResponse.data], {
                                    type: innerResponse.headers['content-type']
                                });
                                console.log(blob.size)
                                zip.file(fileName, blob)
                            }
                        )
                    }
                }
                /*   let mimeType = result.ContentType
                   let fileName = fileToDownload.key
                   let blob = new Blob([result.Body], {type: mimeType})
           
                   photoZip.file(fileName[1], blob)
           */
            })
    }

    const DownloadButtonHandler = async () => {
        if (!userStore.isAuthorized) {
            Store.addNotification({
                title: "Error",
                message: "Перед загрузкой авторизуйтесь",
                type: "warning", // 'default', 'success', 'info', 'warning'
                animationIn: ["animated", "fadeIn"], // animate.css classes that's applied
                animationOut: ["animated", "fadeOut"], // animate.css classes that's applied
                container: "top-right", // where to position the notifications
                dismiss: {
                    duration: 3000
                }
            });
            return
        }
        if (loading || items?.length == 0) {
            return
        }
        if (items == null) {
            return
        }
        setLoading(true)
        setLoadState(0)
        for (let i = 0; i < items.length; i = i + 10) {
            let batch: string[] = []
            if (i + 10 < items.length) {
                batch = items.slice(i, i + 10)
            } else {
                batch = items.slice(i, items.length)
            }
            let zip = new JSZip();
            for (let g = 0; g < batch.length; g = g + 1) {
                try {
                    await DownloadFileFromS3(batch[g], `${(g + 1).toString()}.jpg`, zip)
                    setLoadState((loadState) => loadState + 1)
                } catch {
                    continue
                }
            }
            let count = 0
            zip.forEach(function () {
                count++
            });
            if (count != 0) {
                await zip.generateAsync({ type: "blob" })
                    .then(async function (content: any) {
                        await saveAs(content, `photos${i / 10 + 1}.zip`);
                    });
                DeleteFromCart(batch)
            }
        }
        localStorage.removeItem("SelectedPhotos")
        setItems([])
        setLoading(false)
    }


    const DeleteFromCart = (UIDs: string[]) => {
        let tmp = items
        UIDs.forEach(uid => {
            if (tmp != null) {
                tmp = tmp.filter(function (value) {
                    return value != uid;
                });
            }
        })
        var filtered = tmp
        setItems(filtered)
        localStorage.setItem("SelectedPhotos", JSON.stringify(filtered))

    }

    useEffect(() => { console.log(loadState) }, [loadState])

    return {
        items, DownloadButtonHandler, DeleteFromCart, loading, loadState
    }

}