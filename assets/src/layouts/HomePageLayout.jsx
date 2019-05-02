import React from "react";

export default function MainPageLayout() {
  return (
    <div className="project">
      <div className="title">
        <h4>
          <p>Názov diplomovej práce:</p>
          <span>Verifikačný nástroj pre formalizmus UNITY</span>
          <span className="small">Verification tool for UNITY formalism</span>
          <hr />
          <p>Vedúci práce: </p>
          <span className="small">doc. RNDr. Damas Gruska, PhD.</span>
          <hr />
          <p>Author práce: </p>
          <span className="small">Bc. Filip Špaldoň</span>
        </h4>
      </div>
      <p className="goal">
        Cieľom práce je vytvoriť verifikačný nástroj podporujúci zápis a
        verifikáciu programov zapísaných v jazyku UNITY. Nástroj bude využívať
        transformáciu UNITY programov do vlákien pre každé z priradení a takto
        vzniknutý program sa bude verifikovať nástrojom NuSMV.
      </p>
    </div>
  );
}
