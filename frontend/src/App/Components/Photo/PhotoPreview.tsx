import React, { useContext, useEffect, useState } from 'react';
import { useNavigate, useParams, useLocation } from 'react-router-dom';
import "react-placeholder/lib/reactPlaceholder.css";
import { CloseButton, Col, Modal, Row } from "react-bootstrap";
import { PhotoMeta } from "../../utils/types";
import { useApi } from "../../../api/api";
import "./styles.scss"
import { PhotoWatermark } from './PhotoWatermark';
import { Tag } from '../Tag/Tag';
import { Store } from 'react-notifications-component';
import { AppContext } from '../../AppContext';


interface Photo {
    url: string
}

export const PhotoPreview: React.FC<Photo> = ({ url }) => {

    const [isLoaded, setLoaded] = useState(false)
    const [UUID, setUUID] = useState<string>()
    const [display, setDisplay] = useState("display-none")
    const navigate = useNavigate()
    const location = useLocation();

    const onPhotoSelectedHandler = () => {
        if (UUID != null) {
            console.log(UUID)
            let data = localStorage.getItem("SelectedPhotos")
            if (data != null) {
                let items = JSON.parse(data)
                items.push(UUID)
                items = Array.from(new Set(items))
                localStorage.setItem("SelectedPhotos", JSON.stringify(items))
            } else {
                let items = [UUID]
                localStorage.setItem("SelectedPhotos", JSON.stringify(items))
            }
            Store.addNotification({
                title: "OK",
                message: "Выбранное фото добавлено в корзину",
                type: "success", // 'default', 'success', 'info', 'warning'
                animationIn: ["animated", "fadeIn"], // animate.css classes that's applied
                animationOut: ["animated", "fadeOut"], // animate.css classes that's applied
                container: "top-right", // where to position the notifications
                dismiss: {
                    duration: 2000
                }
            });
        } else {
            Store.addNotification({
                title: "ERROR",
                message: "Ошибка добавления фото в корзину",
                type: "warning", // 'default', 'success', 'info', 'warning'
                animationIn: ["animated", "fadeIn"], // animate.css classes that's applied
                animationOut: ["animated", "fadeOut"], // animate.css classes that's applied
                container: "top-right", // where to position the notifications
                dismiss: {
                    duration: 2000
                }
            });
        }
    }

    const handleImageLoaded = () => {
        setLoaded(true)
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
            {!isLoaded && <div style={{ width: "220px", height: "147px", backgroundColor: "#CCCCCC", marginBottom: "14px" }}></div>}
            <div style={!isLoaded ? { display: "none" } : {}}>
                <div className="item content" onMouseEnter={e => { setDisplay("display-button") }} onMouseLeave={e => { setDisplay("display-none") }}>
                    <img src={url} style={{ margin: "0px 0px" }}
                        onLoad={handleImageLoaded} onClick={showPhotoWatermark} />
                    <button className={"photo-button buy-photo " + display} style={{ position: 'absolute', bottom: 0, left: 0 }} onClick={onPhotoSelectedHandler}>Выбрать фото</button>
                </div>
            </div>
        </>
    )
}

export const MyVerticallyCenteredModal: React.FC = (props: any) => {
    const { api } = useApi()
    const [show, setShow] = useState(true)
    const [photoMeta, setPhotoMeta] = useState<PhotoMeta>()
    const [photoTags, setPhotoTags] = useState<string[]>()
    const { UUID } = useParams()
    const [prevUUIDs, setPrev] = useState<string[]>([])
    const [nextUUIDs, setNext] = useState<string[]>([])
    const { photoStore } = useContext(AppContext)
    const navigate = useNavigate()
    const location = useLocation()
    let state = location.state as { backgroundLocation?: Location };

    useEffect(() => {
        if (photoMeta?.competitors != undefined) {
            setPhotoTags(photoMeta.competitors)
        } else {
            setPhotoTags([])
        }
        getUrls()
    }, [photoMeta]);

    useEffect(() => {
        if (UUID != null) {
            downloadMeta(UUID)
        }
    }, [UUID]);

    const downloadMeta = (UUID: string) => {
        var metaURL = "https://photo.marshalone.ru/api/photo/meta"
        api.get(metaURL, {
            params: {
                UUID: UUID
            }
        })
            .then((response) => {
                let photoMeta: PhotoMeta = response.data
                setPhotoMeta(photoMeta)
            })
            .catch(function (error) {
                if (error.response) {
                    console.log(error.response.data)
                } else if (error.request) {
                    console.log(error.request)
                }
                setPhotoMeta(undefined)
            })
    }

    const getUrls = () => {
        if (photoMeta?.RUID && photoStore[photoMeta?.RUID] != undefined) {
            if (photoMeta.competitors != undefined) {
                if (UUID != undefined) {
                    let index = photoStore[photoMeta.RUID].UUIDs.indexOf(UUID)
                    let tmpPrev: string[] = []
                    for (var i = index - 2; i < index; i++) {
                        if (i < 0 || i > photoStore[photoMeta.RUID].UUIDs.length - 1) {
                            continue
                        } else {
                            tmpPrev.push(photoStore[photoMeta.RUID].UUIDs[i])
                        }
                    }
                    setPrev(tmpPrev)
                    let tmpNext: string[] = []
                    for (var i = index + 1; i <= index + 4 - tmpPrev.length; i++) {
                        if (i < 0 || i > photoStore[photoMeta.RUID].UUIDs.length - 1) {
                            continue
                        } else {
                            tmpNext.push(photoStore[photoMeta.RUID].UUIDs[i])
                        }
                    }
                    setNext(tmpNext)
                }
            } else {
                if (UUID != undefined) {
                    let index = photoStore[photoMeta.RUID].UUIDs.indexOf(UUID)
                    let tmpPrev: string[] = []
                    for (var i = index - 3; i < index; i++) {
                        if (i < 0 || i > photoStore[photoMeta.RUID].UUIDs.length - 1) {
                            continue
                        } else {
                            tmpPrev.push(photoStore[photoMeta.RUID].UUIDs[i])
                        }
                    }
                    setPrev(tmpPrev)
                    let tmpNext: string[] = []
                    for (var i = index + 1; i <= index + 6 - tmpPrev.length; i++) {
                        if (i < 0 || i > photoStore[photoMeta.RUID].UUIDs.length - 1) {
                            continue
                        } else {
                            tmpNext.push(photoStore[photoMeta.RUID].UUIDs[i])
                        }
                    }
                    setNext(tmpNext)
                }
            }
        }
    }

    return (
        <Modal
            show={show}
            aria-labelledby="contained-modal-title-vcenter"
            centered
            size="xl"
            scrollable
            fullscreen
            onHide={() => state.backgroundLocation?.toString() != undefined ? navigate(state.backgroundLocation?.pathname) : navigate("..")}
        >
            <Modal.Header>
                {state?.backgroundLocation != undefined && <CloseButton onClick={() => { state.backgroundLocation?.toString() != undefined ? navigate(state.backgroundLocation?.pathname) : navigate("..") }} />}
            </Modal.Header>
            <Modal.Body style={{ textAlign: "center" }}>
                <Row>
                    <Col >
                        <Row>
                            <Col>
                                <PhotoWatermark url={`https://photo.marshalone.ru/api/photo/file/get?type=watermark&UUID=${UUID}`} />
                            </Col>
                        </Row>
                        {!location.pathname.startsWith("/download") &&
                            <Row style={{ width: "60%", margin: "auto" }}>
                                {prevUUIDs != undefined &&
                                    prevUUIDs.length > 0 &&
                                    prevUUIDs.map(uuid => (
                                        <Col style={{ marginTop: 20 }}>
                                            <img src={`https://photo.marshalone.ru/api/photo/file/get?type=resized&UUID=${uuid}`} className={"preview-opacity modalPreview"} onClick={() => { console.log(state?.backgroundLocation); navigate(`../photo/${uuid}`, { state: { backgroundLocation: state?.backgroundLocation } }) }} />
                                        </Col>
                                    ))
                                }
                                <Col style={{ marginTop: 20 }}>
                                    <img src={`https://photo.marshalone.ru/api/photo/file/get?type=resized&UUID=${UUID}`} className={"modalPreview"} onClick={() => { console.log(state?.backgroundLocation); navigate(`../photo/${UUID}`, { state: { backgroundLocation: state?.backgroundLocation } }) }} />
                                </Col>
                                {nextUUIDs != undefined &&
                                    nextUUIDs.length > 0 &&
                                    nextUUIDs.map(uuid => (
                                        <Col style={{ marginTop: 20 }}>
                                            <img src={`https://photo.marshalone.ru/api/photo/file/get?type=resized&UUID=${uuid}`} className={"preview-opacity modalPreview"} onClick={() => { console.log(state?.backgroundLocation); navigate(`../photo/${uuid}`, { state: { backgroundLocation: state?.backgroundLocation } }) }} />
                                        </Col>
                                    ))
                                }
                            </Row>}
                    </Col>
                    {photoTags && photoTags.length > 0 &&
                        <Col xs={3} style={{ display: "inline-block" }} >
                            {photoTags && photoTags.length > 0 &&
                                <h1>Теги</h1>
                            }
                            {photoTags && photoTags.length > 0 &&
                                photoTags.map((tag) => (
                                    <Tag tag={tag} />))
                            }
                        </Col>
                    }
                </Row>

            </Modal.Body >
        </Modal >
    );
}


// <button className={"photo-button buy-photo " + display} style={{ position: 'absolute', bottom: 0, left: 0 }}>Добавить в корзину</button>



/*const MyVerticallyCenteredModal = (props: any) => {
    return (
        <Modal
            {...props}
            aria-labelledby="contained-modal-title-vcenter"
            centered
            fullscreen
            scrollable
        >
            <Modal.Header closeButton>
            </Modal.Header>
            <Modal.Body style={{ textAlign: "center" }}>
                <Row>
                    <Col >
                        <PhotoWatermark url={watermarkUrl} />
                    </Col>
                    {photoTags && photoTags.length > 0 &&
                        <Col xs={3} style={{ display: "inline-block" }} >
                            {photoTags && photoTags.length > 0 &&
                                <h1>Теги:</h1>
                            }
                            {photoTags && photoTags.length > 0 &&
                                photoTags.map((tag) => (
                                    <Tag tag={tag} />))
                            }
                        </Col>
                    }
                </Row>
            </Modal.Body>
        </Modal >
    );
}*/