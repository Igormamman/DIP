import React, { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import queryString, { parse } from 'query-string'

interface URL {
  url: string
}


export const PhotoWatermark: React.FC<URL> = ({ url }) => {


  return (
    <>
      <img src={url} className="img-fluid"/>
    </>
  )
}

