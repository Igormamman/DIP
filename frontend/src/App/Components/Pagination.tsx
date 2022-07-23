import React, { useEffect } from "react";
import PropTypes from "prop-types";
import { Pagination } from "react-bootstrap";

interface Props {
    currentPage: number;
    totalPages: number;
    handlePage: (event: any) => void;
}
export const MyPagination: React.FC<Props> = ({
    currentPage,
    totalPages,
    handlePage,
}) => {
    var pageArray: (string | number)[] = []
    if (totalPages > 1) {
        if (totalPages <= 9) {
            var i = 1;
            while (i <= totalPages) {
                pageArray.push(i);
                i++;
            }
        } else {
            if (currentPage <= 5) pageArray = [1, 2, 3, 4, 5, 6, 7, 8, "", totalPages];
            else if (totalPages - currentPage <= 4)
                pageArray = [
                    1,
                    "",
                    totalPages - 7,
                    totalPages - 6,
                    totalPages - 5,
                    totalPages - 4,
                    totalPages - 3,
                    totalPages - 2,
                    totalPages - 1,
                    totalPages
                ];
            else
                pageArray = [
                    1,
                    "",
                    currentPage - 3,
                    currentPage - 2,
                    currentPage - 1,
                    currentPage,
                    currentPage + 1,
                    currentPage + 2,
                    currentPage + 3,
                    "",
                    totalPages
                ];
        }
    }

    return (
        <>
            <Pagination style={{ justifyContent: "center" }}>
                {pageArray.map((ele, ind) => {
                    const toReturn = [];
                    if (ind === 0) {
                        toReturn.push(
                            <Pagination.Prev
                                key={"prevpage"}
                                onClick={
                                    currentPage === 1
                                        ? () => { }
                                        : () => {
                                            handlePage(currentPage - 1)
                                        }
                                }
                            />
                        );
                    }
                    if (ele === "") toReturn.push(<Pagination.Ellipsis key={ind}
                        onClick={
                            ind > 2
                                ? () => { handlePage(currentPage + 4) }
                                : () => {
                                    handlePage(currentPage - 4)
                                }
                        }
                    />);
                    else
                        toReturn.push(
                            <Pagination.Item
                                key={ind}
                                active={currentPage === ele ? true : false}
                                onClick={
                                    currentPage === ele
                                        ? () => { }
                                        : () => {
                                            handlePage(ele)
                                        }
                                }
                            >
                                {ele}
                            </Pagination.Item>
                        );
                    if (ind === pageArray.length - 1) {
                        toReturn.push(
                            <Pagination.Next
                                key={"nextpage"}
                                onClick={
                                    currentPage === totalPages - 1
                                        ? () => { }
                                        : () => {
                                            handlePage(currentPage + 1)
                                        }
                                }
                            />
                        );
                    }
                    return toReturn;
                })}
            </Pagination>
        </>
    )

}

