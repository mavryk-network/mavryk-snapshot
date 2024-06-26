import React from 'react'
import { Snapshot } from '../models/Snapshot'

export default function DateSeparator (props: { snapshot: Snapshot; }) {
  return (
    <div style={{
      marginTop: '30px',
      color: '#5F58FF',
      marginBottom: '10px',
      textAlign: 'right',
      fontSize: '2.0vh',
      fontStyle: 'italic'
    }}> {new Date(props.snapshot.date).toISOString().split('T')[0]}
    </div>
  )
}
