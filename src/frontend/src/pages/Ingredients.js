import React from 'react';
import { useForm } from 'react-hook-form';

const About = () => {
  const { register, handleSubmit, formState: { errors } } = useForm();
  const onSubmit = queryIngredients => console.log(queryIngredients); // Query adalah variabel yang bentuknya menyerupai JSON
  console.log(errors);

  return (
          <div className="Ingredients-container">
            <h1 className="Ingredients-title">LIHAT&nbsp;&nbsp;BAHAN</h1>

            <form onSubmit={handleSubmit(onSubmit)}>

                <div className="Search-form">
                  <div className="Search-form-card">
                    <h3>Nama Bahan</h3>
                    <input
                      type="text"
                      placeholder=" "
                      autoComplete='off'
                      className="custom-search-input"
                      {...register("Nama Bahan", { required: true })}
                    />
                  </div>

                  <input 
                  type="submit" 
                  className="submit-button" />
                </div>
            </form>
          
          </div>
        );
};

export default About;
