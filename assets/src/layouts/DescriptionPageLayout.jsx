import React from 'react';

export default function DescriptionPageLayout() {
  return (
    <div className="description">
      <h2>Čo to je UNITY?</h2>
      <p>
        UNITY vychádza z knihy Parallel Program Design - A Foundation, v ktorej
        bol UNITY popísaný a navrhnutý autormi K. Mali Chandy a Jayadev Misra z
        Univerzity of Texax. Je to teoretický jazyk, ktorý sa zameriava na to
        <strong> čo</strong>, namiesto toho <strong> kde</strong>,
        <strong> kedy</strong> alebo<strong> ako</strong>. Jazyk neobsahuje
        žiadnu metódu <strong> riadenia toku</strong> a príkazy programu
        prebiehajú
        <strong> nedeterministickým spôsobom</strong>
        <strong> synchrónne a asynchrónne</strong>, kým sa
        <strong> priradenia</strong> nedostanú do konečného
        <strong> stavu</strong>. To umožňuje, aby programy bežali na neurčito,
        ako napríklad autopilot alebo elektráreň, ktoré by normálne skončili.
      </p>
    </div>
  );
}
