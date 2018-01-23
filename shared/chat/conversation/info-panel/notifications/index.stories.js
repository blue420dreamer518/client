// @flow
import React from 'react'
import {storiesOf, action} from '../../../../stories/storybook'
import {Notifications} from './index'

const load = () => {
  storiesOf('Chat/Conversation/InfoPanelNotifications', module)
    .add('Notifications (unsaved)', () => (
      <Notifications
        channelWide={false}
        desktop="atmention"
        mobile="never"
        resetSaveState={() => {}}
        saveState="unsaved"
        onSetDesktop={action('onSetDesktop')}
        onSetMobile={action('onSetMobile')}
        onToggleChannelWide={action('onToggleChannelwide')}
      />
    ))
    .add('Notifications (saving)', () => (
      <Notifications
        channelWide={true}
        desktop="generic"
        mobile="atmention"
        resetSaveState={() => {}}
        saveState="saving"
        onSetDesktop={action('onSetDesktop')}
        onSetMobile={action('onSetMobile')}
        onToggleChannelWide={action('onToggleChannelwide')}
      />
    ))
    .add('Notifications (saved)', () => (
      <Notifications
        channelWide={true}
        desktop="generic"
        mobile="atmention"
        resetSaveState={() => {}}
        saveState="saved"
        onSetDesktop={action('onSetDesktop')}
        onSetMobile={action('onSetMobile')}
        onToggleChannelWide={action('onToggleChannelwide')}
      />
    ))
}

export default load
