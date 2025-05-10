import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import '../App.css';
import ingredientsData from '../data/ingredients.json';

const About = () => {
  const { register, handleSubmit } = useForm();
  const [searchResults, setSearchResults] = useState([]);
  const [searchPerformed, setSearchPerformed] = useState(false);

  const onSubmit = ({ nama }) => {
    const result = ingredientsData.filter(ingredient =>
      ingredient.name.toLowerCase().includes(nama.toLowerCase())
    );
    setSearchResults(result);
    setSearchPerformed(true);
  };

  const getImagePath = (name) => {
    return require(`../media/ingredients/${name}.svg`);
};


  return (
    <div className="Ingredients-container">
      <h1 className="Ingredients-title">LIHAT&nbsp;&nbsp;&nbsp;BAHAN</h1>

      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="Search-form">
          <div className="Search-form-card">
            <h3>Nama Bahan</h3>
            <input
              type="text"
              placeholder="Contoh: Babe the blue ox"
              autoComplete="off"
              className="custom-search-input"

              {...register('nama', { required: true })}
            />
          </div>
          {/* <input type="submit" className="submit-button" value="Cari" /> */}
        </div>
      </form>

      {searchPerformed && (
        <>
          <h2 className="Ingredients-subtitle">HASIL PENCARIAN</h2>
          {searchResults.length === 0 ? (
            <p className="Ingredients-noresult">Tidak ditemukan.</p>
          ) : (
            <div className="Ingredient-grid">
              {searchResults.map((item, index) => (
                <div key={index} className="Ingredient-card">
                  <img src={getImagePath(item.name)} alt={item.name} className="Ingredient-icon" />
                  <p>{item.name}</p>
                </div>
              ))}
            </div>
          )}
        </>
      )}

      <h2 className="Ingredients-subtitle">LIHAT SEMUA</h2>
      <div className="Ingredient-grid">
        {ingredientsData.map((item, index) => (
          <div key={index} className="Ingredient-card">
            <img src={getImagePath(item.name)} alt={item.name} className="Ingredient-icon" />
            <p>{item.name}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default About;