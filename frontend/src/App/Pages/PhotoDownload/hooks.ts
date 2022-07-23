import React, { useEffect, useState } from "react";
import { useApi } from "../../../api/api"
import { useNavigate, useParams } from "react-router-dom";
import watermark from "../../utils/images/Gerb_mgtu.png"
import { Previews, Race, WidthHeight } from "../../utils/types";
import { Store } from "react-notifications-component";
import { ServerRouts } from "../../utils/routs";
import moment from "moment";
import { debounce } from "../../utils/common";

export const useDownloadPage = () => {

    const { api } = useApi()
    const [selectedFile, setSelectedFile] = useState<File[]>();
    const [fileCount, setFileCount] = useState<number>(1);
    const [currentFileNum, setCurrentFileNum] = useState<number>(0);
    const [UUIDs, setUUIDs] = useState<string[]>([])
    const [loading, setLoading] = useState(false)
    const image = new Image()
    const navigate = useNavigate()
    const [access, setAccess] = useState<boolean>()
    const [page, setPage] = useState<number>(1);
    const [totalPages, setTotalPages] = useState(1);
    const [raceAccess, setRaceAccess] = useState<Race>()
    const { raceUID } = useParams();

    const handlePageClick = (event: any) => {
        setPage(event);
    };

    const getRaceAccess = async () => {
        api.get(ServerRouts.GetRaceAccess, {
            params: {
                ruid: raceUID
            }
        }).then((response) => {
            let data: Race
            if (response.data != undefined) {
                data = response.data
                if (data != null) {
                    setRaceAccess(data);
                }
            }
        })
    }


    const getRaceInfo = async () => {
        api.get(ServerRouts.GetRaceInfo, {
            params: {
                ruid: raceUID
            }
        }).then((response) => {
            let data: Race
            if (response.data != undefined) {
                data = response.data
                if (data != null) {
                    if (moment().diff(data.date, 'days') >= 0) {
                        setAccess(true)
                    } else {
                        setAccess(false)
                    }
                }
            }
        })
    }

    useEffect(() => {
        getRaceInfo()
        getRaceAccess()
    }, [])

    interface NewFile {
        token: string;
        fileType: string;
    }

    interface FileMeta {
        photoName: string;
        imageSize: Number;
        resizedSize: Number;
        mlSize: Number;
        watermarkSize: Number;
        imageWidth: Number;
        imageHeight: Number;
    }

    interface FileData {
        fileName: string
    }

    interface Task {
        fileData: FileData[]
        RaceID: string
    }

    const resizeFile = (file: File) => {
        return new Promise(resolve => {
            var height: number
            var width: number
            var contentType = "image/jpeg";
            var canvas = document.createElement('canvas');
            var ctx = canvas.getContext('2d');
            const reader = new FileReader();
            reader.onload = () => {
                var image = new Image();
                image.onload = function () {
                    if (image.height > image.width) {
                        height = 480
                        width = 480 * (image.width / image.height)
                    } else {
                        height = 480 * (image.height / image.width)
                        width = 480
                    }
                    canvas.height = height
                    canvas.width = width
                    if (ctx != null) {
                        ctx.drawImage(image, 0, 0, canvas.width, canvas.height);
                        canvas.toBlob(function (blob) {
                            if (blob != null) {
                                var newFile = new File([blob], file.name, {
                                    type: contentType,
                                    lastModified: Math.floor(new Date().getTime() / 1000)
                                })
                                resolve(newFile)
                            }
                        }, 'image/jpeg', 1)
                    }
                }
                if (reader.result?.toString() != null) {
                    image.src = reader.result?.toString();
                }
            }
            reader.readAsDataURL(file);
        })
    }

    const resizeToML = (file: File) => {
        return new Promise(resolve => {
            var height: number
            var width: number
            var contentType = "image/jpeg";
            var canvas = document.createElement('canvas');
            var ctx = canvas.getContext('2d');
            const reader = new FileReader();
            reader.onload = () => {
                var image = new Image();
                image.onload = () => {
                    if (image.height > image.width) {
                        height = 1536
                        width = 1536 * (image.width / image.height)
                    } else {
                        height = 1536 * (image.height / image.width)
                        width = 1536
                    }
                    canvas.height = height
                    canvas.width = width
                    ctx?.drawImage(image, 0, 0, canvas.width, canvas.height)
                    canvas.toBlob(function (blob) {
                        if (blob != null) {
                            var newFile = new File([blob], file.name, {
                                type: contentType,
                                lastModified: Math.floor(new Date().getTime() / 1000)
                            })
                            resolve(newFile)
                        }
                    }, 'image/jpeg', 1)

                }
                if (reader.result?.toString() != null) {
                    image.src = reader.result?.toString();
                }
            }
            reader.readAsDataURL(file);
        })
    }

    const addWatermark = (file: File) => {
        return new Promise(resolve => {
            var height: number
            var width: number
            var contentType = "image/jpeg";
            var canvas = document.createElement('canvas');
            var ctx = canvas.getContext('2d');
            const reader = new FileReader();
            reader.onload = () => {
                var image = new Image();
                image.onload = () => {
                    if (image.height > image.width) {
                        height = 1024
                        width = 1024 * (image.width / image.height)
                    } else {
                        height = 1024 * (image.height / image.width)
                        width = 1024
                    }
                    canvas.height = height
                    canvas.width = width
                    ctx?.drawImage(image, 0, 0, canvas.width, canvas.height)
                    var watermarkImage = new Image()
                    watermarkImage.src = watermark
                    watermarkImage.onload = () => {
                        if (ctx != null) {
                            ctx.globalAlpha = 1;
                            ctx?.drawImage(watermarkImage, 0, 0, width / 5, height / 5)
                        }
                        canvas.toBlob(function (blob) {
                            if (blob != null) {
                                var newFile = new File([blob], file.name, {
                                    type: contentType,
                                    lastModified: Math.floor(new Date().getTime() / 1000)
                                })
                                resolve(newFile)
                            }
                        }, 'image/jpeg', 1)
                    }
                }
                if (reader.result?.toString() != null) {
                    image.src = reader.result?.toString();
                }
            }
            reader.readAsDataURL(file);
        })
    }

    const getWidthHeight = (file: File) => {
        return new Promise(resolve => {
            const reader = new FileReader();
            var image = new Image();
            reader.onload = () => {
                image.onload = () => {
                    var widthHeigth: WidthHeight = { width: image.width, height: image.height }
                    resolve(widthHeigth)
                }
                if (reader.result?.toString() != null) {
                    image.src = reader.result?.toString();
                }
            }
            reader.readAsDataURL(file);
        })
    }

    const fileSelectedHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
        const fileList = event.target.files;
        if (!fileList) {
            return;
        }
        let files: File[] = []
        for (var i = 0; i < fileList.length; i++) {
            files.push(fileList[i])
        }
        setSelectedFile(files)
    }

    const uploadFileHandler = async function (e: React.MouseEvent<HTMLSpanElement, MouseEvent>) {

        if (loading || selectedFile == undefined) {
            return
        }

        setLoading(true)
        if (!selectedFile) {
            return
        }
        setCurrentFileNum(0)
        setFileCount(selectedFile.length)
        console.log(selectedFile.length)
        /*for (let i = 0; i < selectedFile.length; i++) {
            sendFile(selectedFile[i])
        }*/
        sendTask(selectedFile).then(async () => {
            for (let i = 0; i < selectedFile.length; i = i + 1) {
                setCurrentFileNum(i)
                await sendFile(selectedFile[i])
            }
            setLoading(false)
        }, () => {
            setLoading(false)
        })
        /*  {
    
              api.post("http://localhost:4000/photo", formData)
                  .then((response) => {
                      if (response.status == 200) {
                          console.log("молодец", response)
                      }
                  })
                  .catch((error) => [
                      console.log(error)
                  ])
          }*/
    }
    //  api.post("https://api-dip.duckdns.org/photo", formData)
    //  api.post("http://localhost:4000/photo", formData)

    const sendFileBatch = async (fileList: File[]) => {
        await Promise.all(fileList.map(file => sendFile(file)))
    }

    const selectedFileInfo = () => {
        if (selectedFile?.length == undefined || selectedFile?.length == 0) {
            return "Выберите файлы для загрузки"
        } else if (selectedFile?.length == 1) {
            return `Выбран 1 файл, ${selectedFile[0].name}`
        } else if (selectedFile?.length < 5) {
            return `Выбрано ${selectedFile?.length} файла`
        } else {
            return `Выбрано ${selectedFile?.length} файлов`
        }
    }

    const loadFile = (fileInfo: NewFile, file: File) => {
        return new Promise<void>(resolve => {
            var fileForm = new FormData();
            fileForm.set('token', fileInfo.token)
            fileForm.set('file_type', fileInfo.fileType)
            fileForm.set('file', file)
            api.post("https://photo.marshalone.ru/api/photo/file/upload", fileForm)
                .then((response) => {
                    if (response.status == 200) {
                        resolve()
                    }
                })
                .catch((error) => {
                    console.log("ASDASDASD")
                    setLoading(false)
                    resolve();
                    throw (error);
                })
        })
    }
    //   api.post("https://api-dip.duckdns.org/file/upload", fileForm)
    //  api.post("http://localhost:4004/file/upload", fileForm)
    //  api.post("http://photo.marshalone.ru:4004/file/upload", fileForm)

    const deleteHandler = () => {
        getPreviews()
    }

    const getPreviews = () => {
        api.get(ServerRouts.Previews, {
            params: {
                raceUID: raceUID,
                limit: 24,
                offset: 24*(page-1),
                photographer: "true",
            }
        }).then((response) => {
            let data: Previews
            if (response.data != undefined) {
                data = response.data
            } else {
                data = { count: 0, previewURL: [] }
            }
            setUUIDs(data.previewURL)
            if (data.count % 24 == 0) {
                setTotalPages(Math.floor(data.count / 24))
            } else {
                setTotalPages(Math.floor(data.count / 24) + 1)
            }
        })
    }

    const sendTask = (files: File[]) => {
        return new Promise<void>(resolve => {
            if (raceUID != undefined) {
                var task: Task = { RaceID: raceUID, fileData: [] }
                if (!files) {
                    return
                }
                for (let i = 0; i < files.length; i = i + 1) {
                    task.fileData.push({ fileName: files[i].name })
                }
                api.post(ServerRouts.NewTask, task).then((response) => {
                    if (response.status == 200) {

                    } else {
                        Store.addNotification({
                            title: "Error",
                            message: "Ошибка загрузки (у вас точно есть доступ??)",
                            type: "warning", // 'default', 'success', 'info', 'warning'
                            animationIn: ["animated", "fadeIn"], // animate.css classes that's applied
                            animationOut: ["animated", "fadeOut"], // animate.css classes that's applied
                            container: "top-right", // where to position the notifications
                            dismiss: {
                                duration: 3000
                            }
                        });
                        setLoading(false)
                    }
                    resolve();
                })
            }
        })
    }

    const sendFile = (file: File) => {
        return new Promise<void>(resolve => {
            resizeFile(file).then((resizedFile: any) => {
                resizeToML(file).then((resizedToMLFile: any) => {
                    addWatermark(resizedToMLFile).then((watermarkFile: any) => {
                        getWidthHeight(file).then((widthHeight: any) => {
                            var fileMeta: FileMeta = {
                                photoName: file.name,
                                imageSize: file.size,
                                resizedSize: resizedFile.size,
                                mlSize: resizedToMLFile.size,
                                watermarkSize: watermarkFile.size,
                                imageWidth: widthHeight.width,
                                imageHeight: widthHeight.height
                            };
                            api.post(ServerRouts.MetaUpdate, fileMeta)
                                .then((response) => {
                                    if (response.status == 200) {
                                        console.log("молодец", response)
                                        var fileInfo: NewFile[] = response.data
                                        try {
                                            //  resolve()
                                            loadFile(fileInfo[0], file).then(() => loadFile(fileInfo[1], resizedFile))
                                                .then(() => loadFile(fileInfo[2], watermarkFile))
                                                .then(() => loadFile(fileInfo[3], resizedToMLFile)).then(() => resolve())
                                            Store.addNotification({
                                                title: "Success",
                                                message: "Файл успешно загружен " + file.name,
                                                type: "success", // 'default', 'success', 'info', 'warning'
                                                animationIn: ["animated", "fadeIn"], // animate.css classes that's applied
                                                animationOut: ["animated", "fadeOut"], // animate.css classes that's applied
                                                container: "top-right", // where to position the notifications
                                                dismiss: {
                                                    duration: 3000
                                                }
                                            });
                                        } catch (e: unknown) {
                                            if (e instanceof Error) {
                                                Store.addNotification({
                                                    title: "Error",
                                                    message: "Файл не загружен " + file.name,
                                                    type: "default", // 'default', 'success', 'info', 'warning'
                                                    animationIn: ["animated", "fadeIn"], // animate.css classes that's applied
                                                    animationOut: ["animated", "fadeOut"], // animate.css classes that's applied
                                                    container: "top-right", // where to position the notifications
                                                    dismiss: {
                                                        duration: 3000
                                                    }
                                                });
                                            }
                                        }
                                    }
                                })
                                .catch((e) => {
                                    Store.addNotification({
                                        title: "Error",
                                        message: "Файл не загружен " + file.name,
                                        type: "default", // 'default', 'success', 'info', 'warning'
                                        animationIn: ["animated", "fadeIn"], // animate.css classes that's applied
                                        animationOut: ["animated", "fadeOut"], // animate.css classes that's applied
                                        container: "top-right", // where to position the notifications
                                        dismiss: {
                                            duration: 3000
                                        }
                                    });

                                    console.log(e)
                                })
                        })
                    })
                })
            })
        })
    }

    useEffect(() => {
        image.src = watermark
        getPreviews()
    }, []);

    useEffect(() => {
        if (loading == false) {
            setCurrentFileNum(0)
            setFileCount(1)
        }
    }, [loading]);

    const debouncedGetPreviews = debounce(getPreviews, 500);

    useEffect(() => {
            debouncedGetPreviews()
    }, [page]);

    // api.get("http://localhost:4000/races")
    // api.get("https://api-dip.

    const gotoGallery = () => {
        navigate("../")
    }

    return {
        fileSelectedHandler,
        uploadFileHandler,
        access,
        gotoGallery,
        selectedFileInfo,
        loading,
        fileCount,
        currentFileNum,
        handlePageClick,
        page,
        totalPages,
        UUIDs,
        deleteHandler,
        raceAccess
    }

}