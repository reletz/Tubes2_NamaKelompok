import React from 'react';
import nau from '../media/nau.png';
import lucas from '../media/lucas.png';
import lana from '../media/lana.png';
import '../App.css';

const About = () => {
  return (
    <>
      {/* <div className="About-container"> */}
        <div className="About-title">
          <p>TENTANG&nbsp;&nbsp;&nbsp;KAMI</p>
        </div>
        <div className="About-image-avatars">
          <img src={lana} className="Lana-avatar" alt="Lana Avatar"/>
          <img src={lucas} className="Lucas-avatar" alt="Lana Avatar"/>
          <img src={nau} className="Nau-avatar" alt="Lana Avatar"/>
        </div>
        <div className="About-authors-lana">
          <p>SAMANTHA LAQUEENNA GINTING<br/>NIM 13523138</p>
          <p><br/>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam ultrices turpis eu libero elementum, ut cursus turpis commodo. Vestibulum dui arcu, aliquet vitae mauris vel, congue mollis nunc. Etiam quam mi, convallis ac elit id, sagittis cursus eros.</p>
        </div>
        <div className="About-authors-lucas">
          <p>NICHOLAS ANDHIKA LUCAS<br/>NIM 13523014</p>
          <p><br/>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam ultrices turpis eu libero elementum, ut cursus turpis commodo. Vestibulum dui arcu, aliquet vitae mauris vel, congue mollis nunc. Etiam quam mi, convallis ac elit id, sagittis cursus eros.</p>
        </div>
        <div className="About-authors-nau">
          <p>NAUFARREL ZHAFIF ABHISTA<br/>NIM 13523149</p>
          <p><br/>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam ultrices turpis eu libero elementum, ut cursus turpis commodo. Vestibulum dui arcu, aliquet vitae mauris vel, congue mollis nunc. Etiam quam mi, convallis ac elit id, sagittis cursus eros.</p>
        </div>
      {/* </div> */}
    </>
        );
};

export default About;
