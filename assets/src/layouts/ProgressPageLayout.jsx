import React from "react";
import { Timeline, TimelineEvent } from "react-event-timeline";

export default function ProgressPageLayout() {
  return (
    <Timeline className="timeline">
      <TimelineEvent
        title="Editácia prvej kapitoly"
        createdAt="27. januára 2019"
        icon={<i className="material-icons md-18" />}
      >
        Podľa pripomienok som upravil prvú kapitolu.
      </TimelineEvent>
      <TimelineEvent
        title="Implementácia interpretra"
        createdAt="24. januára 2019 - 26. januára 2019"
        icon={<i className="material-icons md-18" />}
      >
        Implemantácia interpretra - čítanie vstupu, vytváranie tokenov, scaner,
        parser.
      </TimelineEvent>
      <TimelineEvent
        title="Implementácia MonacoEditoru"
        createdAt="23. januára 2019 - 24. januára 2019"
        icon={<i className="material-icons md-18" />}
      >
        Implemantácia editoru, nastavovanie potrebného config súboru na
        rozoznávanie syntaxu.
      </TimelineEvent>
      <TimelineEvent
        title="Výber vhodného react componentu na code editor"
        createdAt="22. januára 2019"
        icon={<i className="material-icons md-18" />}
      >
        Vybral som si MonacoEditor react component na prostredie pre úpravu
        vstupného kódu.
      </TimelineEvent>
      <TimelineEvent
        title="Projektový seminár - 1. opravný termín"
        createdAt="21. januára 2019"
        icon={<i className="material-icons md-18" />}
      >
        Na skúške mi chýbala implemantácia, ktorú musím spraviť do 28.1.
      </TimelineEvent>
      <TimelineEvent
        title="www stránka diplomovky - stav online"
        createdAt="20. januára 2019"
        icon={<i className="material-icons md-18" />}
      >
        Svoju stránku na diplomovú prácu som nasadil online.
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
        Na stretunutí sme sa dohodli o kapitole, ktorú budem písať ako prvú a
        ako ukážku na predmet projektový seminár 2.
      </TimelineEvent>
    </Timeline>
  );
}
