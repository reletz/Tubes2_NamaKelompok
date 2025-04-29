import React from 'react';
import nau from '../media/nau.png';
import lucas from '../media/lucas.png';
import lana from '../media/lana.png';
import '../App.css';

const About = () => {
  return (
    <div className="About-container">
      <h1 className="About-title">TENTANG&nbsp;&nbsp;KAMI</h1>
      <div className="About-cards">
        <div className="About-card">
          <img src={lana} className="About-avatar" alt="Lana Avatar" />
          <div className="About-info">
            <h2>SAMANTHA LAQUEENNA GINTING</h2>
            <h3>NIM 13523138</h3>
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam ultrices turpis eu libero elementum, ut cursus turpis commodo. Vestibulum dui arcu, aliquet vitae mauris vel, congue mollis nunc. Etiam quam mi, convallis ac elit id, sagittis cursus eros.</p>
          </div>
        </div>
        <div className="About-card">
          <img src={lucas} className="About-avatar" alt="Lucas Avatar" />
          <div className="About-info">
            <h2>NICHOLAS ANDHIKA LUCAS</h2>
            <h3>NIM 13523014</h3>
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam ultrices turpis eu libero elementum, ut cursus turpis commodo. Vestibulum dui arcu, aliquet vitae mauris vel, congue mollis nunc. Etiam quam mi, convallis ac elit id, sagittis cursus eros.</p>
          </div>
        </div>
        <div className="About-card">
          <img src={nau} className="About-avatar" alt="Nau Avatar" />
          <div className="About-info">
            <h2>NAUFARREL ZHAFIF ABHISTA</h2>
            <h3>NIM 13523149</h3>
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam ultrices turpis eu libero elementum, ut cursus turpis commodo. Vestibulum dui arcu, aliquet vitae mauris vel, congue mollis nunc. Etiam quam mi, convallis ac elit id, sagittis cursus eros.</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default About;
