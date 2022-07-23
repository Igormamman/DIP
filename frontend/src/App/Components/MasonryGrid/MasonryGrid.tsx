import React, { useEffect } from 'react';
import "./styles.scss"
import Masonry from "react-masonry-css";
import {PhotoPreview} from "../Photo/PhotoPreview";

const Columns4 = {
    default: 4,
    1000: 3,
    750: 2,
    500: 1
};
const Columns3 = {
    default: 3,
    750: 2,
    500: 1
};
const Columns2 = {
    default: 2,
    500: 1,
};
const Columns1 = {
    default: 1,
};


interface Photos {
    uuids: string[]
}


export const PhotoGrid: React.FC<Photos> = ({uuids}) => {

    if (uuids == null || uuids.length == 0) {
        return (<>
            <div style={{margin: "0 , auto", width: "fit-content", maxWidth: "1000px", minHeight: "50px"}}>
            </div>
        </>)
    }

    switch (uuids.length) {
        case 1: {
            return (
                <>
                    <div style={{margin: "auto", width: "fit-content", maxWidth: "1000px"}}>
                        <Masonry breakpointCols={Columns1}
                                 className="my-masonry-grid"
                                 columnClassName="my-masonry-grid_column">
                            {uuids.length>0 &&
                            uuids.map((uuid) => (
                                <PhotoPreview url={`https://photo.marshalone.ru/api/photo/file/get?UUID=${uuid}&type=resized`}/>))
                            }
                        </Masonry>
                    </div>
                </>
            )
        }
        case 2: {
            return (
                <>
                    <div style={{margin: "auto", width: "fit-content", maxWidth: "1000px"}}>
                        <Masonry breakpointCols={Columns2}
                                 className="my-masonry-grid"
                                 columnClassName="my-masonry-grid_column">
                            {uuids &&
                            uuids.map((uuid) => (
                                <PhotoPreview url={`https://photo.marshalone.ru/api/photo/file/get?UUID=${uuid}&type=resized`}/>))
                            }
                        </Masonry>
                    </div>
                </>
            )
        }
        case 3: {
            return (
                <>
                    <div style={{margin: "auto", width: "fit-content", maxWidth: "1000px"}}>
                        <Masonry breakpointCols={Columns3}
                                 className="my-masonry-grid"
                                 columnClassName="my-masonry-grid_column">
                            {uuids &&
                            uuids.map((uuid) => (
                                <PhotoPreview url={`https://photo.marshalone.ru/api/photo/file/get?UUID=${uuid}&type=resized`}/>))
                            }
                        </Masonry>
                    </div>
                </>
            )
        }
        case 4: {
            return (
                <>
                    <div style={{margin: "auto", width: "fit-content", maxWidth: "1000px"}}>
                        <Masonry breakpointCols={Columns4}
                                 className="my-masonry-grid"
                                 columnClassName="my-masonry-grid_column">
                            {uuids &&
                            uuids.map((uuid) => (
                                <PhotoPreview url={`https://photo.marshalone.ru/api/photo/file/get?UUID=${uuid}&type=resized`}/>))
                            }
                        </Masonry>
                    </div>
                </>
            )
        }
        default: {
            return (
                <>
                    <div style={{margin: "auto", width: "fit-content", maxWidth: "1000px"}}>
                        <Masonry breakpointCols={Columns4}
                                 className="my-masonry-grid"
                                 columnClassName="my-masonry-grid_column">
                            {uuids &&
                            uuids.map((uuid) => (
                                <PhotoPreview url={`https://photo.marshalone.ru/api/photo/file/get?UUID=${uuid}&type=resized`}/>))
                            }
                        </Masonry>
                    </div>
                </>
            )
        }

    }


}