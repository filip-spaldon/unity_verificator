import React from 'react';

export default function MainPageLayout() {
  return (
    <div className="project">
      <div className="title">
        <h4>
          <p>Názov diplomovej práce:</p>
          <span>Verifikačný nástroj pre UNITY</span>
          <span className="small">Verification tool for UNITY</span>
          <hr />
          <p>Vedúci práce: </p>
          <span className="small">doc. RNDr. Damas Gruska, PhD.</span>
          <hr />
          <p>Author práce: </p>
          <span className="small">Bc. Filip Špaldoň</span>
        </h4>
      </div>
      <p className="goal">
        Cieľom práce je vytvoriť verifikačný nástroj podporujúci zápis,
        simuláciu a verifikáciu programov zapísaných v jazyku UNITY.
      </p>
    </div>
  );
}
