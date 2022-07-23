import {useGalleryPage} from "./hooks";
import {PhotoWatermark} from '../../Components/Photo/PhotoWatermark';
import {Header} from "../../Components/M1Header/Header";

export const WatermarkPage: React.FC = () => {

    const {
        url,
        gotoGallery,
        meta
    } = useGalleryPage()

    return (
        <div>
            <Header/>
            <div style={{marginTop: 40, textAlign: "center"}}>
                <PhotoWatermark url={url}/>
                {meta?.competitors &&
                <h2>{meta.competitors}</h2>}
                <button className='btn btn-secondary' style={{marginTop: 20}} onClick={gotoGallery}>Назад</button>
            </div>
        </div>
    )
};

