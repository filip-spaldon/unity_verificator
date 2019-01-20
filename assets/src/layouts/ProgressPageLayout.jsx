import React from 'react';
import { Timeline, TimelineEvent } from 'react-event-timeline';

export default function ProgressPageLayout() {
  return (
    <Timeline>
      <TimelineEvent
        title="www stránka diplomovky - stav online"
        createdAt="20. januára 2019"
        icon={<i className="material-icons md-18" />}
      >
        Svoje stránku na diplomovú prácu som nasadil online.
      </TimelineEvent>
      <TimelineEvent
        title="1 kapitola práce"
        createdAt="18. januára 2019"
        icon={<i className="material-icons md-18" />}
      >
        Dokončenie prvej kapitoly o dĺžke 16 strán z diplomovej práce.
      </TimelineEvent>
      <TimelineEvent
        title="Stretnutie so školiteľom"
        createdAt="13. decembra 2018"
        icon={<i className="material-icons md-18" />}
      >
        Na stretunutí sme sa dohodli o kapitole, ktorú budem písať ako prvá a
        ako ukážku na predmet projektový seminár 2.
      </TimelineEvent>
    </Timeline>
  );
}
