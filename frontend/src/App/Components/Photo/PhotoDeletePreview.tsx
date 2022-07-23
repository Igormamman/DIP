import React, { useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import "react-placeholder/lib/reactPlaceholder.css";
import { useApi } from "../../../api/api";
import "./styles.scss"
import { Store } from 'react-notifications-component';
import { ServerRouts } from '../../utils/routs';


interface Photo {
    url: string
    deleteHandler: Function
}

export const PhotoDeletePreview: React.FC<Photo> = ({ url, deleteHandler }) => {

    const [isLoaded, setLoaded] = useState(false)
    const { api } = useApi()
    const [UUID, setUUID] = useState<string>()
    const [display, setDisplay] = useState("display-none")
    const navigate = useNavigate();
    const location = useLocation();

    const handleImageLoaded = () => {
        setLoaded(true)
    }


    const deletePhoto = (url: string) => {
        let myURL = new URL(url);
        if (myURL.searchParams.has('UUID') == true) {
            api.get(ServerRouts.DeletePhoto, {
                params: {
                    UUID: myURL.searchParams.get('UUID'),
                }
            }).then((response) => {
                if (response.status == 200) {
                    Store.addNotification({
                        title: "OK",
                        message: "Фото успешно удалено",
                        type: "success", // 'default', 'success', 'info', 'warning'
                        animationIn: ["animated", "fadeIn"], // animate.css classes that's applied
                        animationOut: ["animated", "fadeOut"], // animate.css classes that's applied
                        container: "top-right", // where to position the notifications
                        dismiss: {
                            duration: 2000
                        }
                    });
                    deleteHandler()
                }
                if (response.status == 400 || response.status == undefined) {
                    Store.addNotification({
                        title: "Error",
                        message: "Ошибка удаления фото",
                        type: "warning", // 'default', 'success', 'info', 'warning'
                        animationIn: ["animated", "fadeIn"], // animate.css classes that's applied
                        animationOut: ["animated", "fadeOut"], // animate.css classes that's applied
                        container: "top-right", // where to position the notifications
                        dismiss: {
                            duration: 2000
                        }
                    });
                    deleteHandler()
                }
            }).catch(() => {
                Store.addNotification({
                    title: "Error",
                    message: "Ошибка удаления фото",
                    type: "warning", // 'default', 'success', 'info', 'warning'
                    animationIn: ["animated", "fadeIn"], // animate.css classes that's applied
                    animationOut: ["animated", "fadeOut"], // animate.css classes that's applied
                    container: "top-right", // where to position the notifications
                    dismiss: {
                        duration: 2000
                    }
                });
                deleteHandler()
            })
        }
    }

    useEffect(() => {
        setLoaded(false)
        let URLobj = new URL(url)
        if (URLobj.searchParams.has('UUID') == true) {
            let uuid = URLobj.searchParams.get('UUID')
            if (uuid != null) {
                setUUID(uuid)
            }
        }
    }, [url]);

    const showPhotoWatermark = () => {
        navigate(`../../photo/${UUID}`, { state: { backgroundLocation: location } })
    }

    return (
        <>
            {!isLoaded && <div className='mx-auto' style={{ width: "480px", height: "320px", backgroundColor: "#CCCCCC", marginBottom: "20px" }}></div>}
            <div className="content photo-center item delete" onMouseEnter={e => { setDisplay("display-button") }} onMouseLeave={e => { setDisplay("display-none") }} style={!isLoaded ? { display: "none" } : {}}>
                <img src={url} style={{ margin: "0px 0px" }}
                    onLoad={handleImageLoaded}/>
                <button className={"photo-button buy-photo " + display} style={{ position: 'absolute', bottom: 0, left: 0 }} onClick={() => deletePhoto(url)}>Удалить фото</button>
            </div>
        </>
    )
}