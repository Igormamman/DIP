import React from 'react'
import { Badge } from 'react-bootstrap'
import "./styles.scss"

interface TagInfo {
  tag: string
}

export const Tag: React.FC<TagInfo> = ({ tag }) => {
  return (
    <>
        <Badge bg="warning" className="tag">{tag}</Badge>
    </>
  )
}

