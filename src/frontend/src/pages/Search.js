import React from 'react';
import { useForm } from 'react-hook-form';
import Tree from 'react-d3-tree';
import productTree from '../data/product_tree.json';

const About = () => {
  const { register, handleSubmit } = useForm();
  const [treeData, setTreeData] = React.useState(null);

  const onSubmit = (querySearch) => {
    console.log(querySearch);
    setTreeData(productTree);
  };

  const treeContainerRef = React.useRef(null);
  const [dimensions, setDimensions] = React.useState({ width: 0, height: 0 });

  React.useEffect(() => {
    if (treeContainerRef.current) {
      const { offsetWidth, offsetHeight } = treeContainerRef.current;
      setDimensions({ width: offsetWidth, height: offsetHeight });
    }
  }, [treeData]);

  const renderCustomNode = ({ nodeDatum }) => {
    const isLeaf = !nodeDatum.children || nodeDatum.children.length === 0;

    return (
      <g>
        <foreignObject width={90} height={55} x={-50} y={0}>
          <div
            xmlns="http://www.w3.org/1999/xhtml"
            className={`tree-node ${isLeaf ? 'leaf-node' : 'parent-node'}`}
          >
            <p>{nodeDatum.name}</p>
          </div>
        </foreignObject>
      </g>
    );
  };

  return (
    <div className="Search-container">
      <h1 className="Search-title">CARI&nbsp;&nbsp;&nbsp;RESEP</h1>
      
      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="Search-form">
          <div className="Search-form-card">
            <h3>Nama Resep*</h3>
            <input
              type="text"
              placeholder="Contoh: Babe the blue ox"
              autoComplete='off'
              className="custom-search-input"
              {...register("Nama Resep", { required: true })}
            />
          </div>

          <div className="Search-form-card">
            <h3>Jumlah Resep*</h3>
            <input 
              type="number"
              placeholder="Contoh: 5"
              className="custom-search-input"
              {...register("Maksimal Resep", {required: true, min: 1})} />
          </div>

          <div className="Search-form-card">
            <h3>Algoritma*</h3>
            <div className="Search-form-card-radio">
              <div className="Search-form-card-radio-2">
                <label className="custom-radio">
                  <input
                    type="radio"
                    {...register("Algoritma", { required: true })}
                    value="BFS"
                  />
                  <span className="radio-image" />
                </label>
                
                <label className="custom-radio">
                  <input
                    type="radio"
                    {...register("Algoritma", { required: true })}
                    value="DFS"
                  />
                  <span className="radio-image" />
                </label>
              </div>
              <div className="Search-form-card-radio-3">
                <p>BFS</p>
                <p>DFS</p>
              </div>
            </div>
          </div>

          <input 
            type="submit" 
            className="submit-button" />
        </div>
      </form>

      {/* Tree View */}
      {treeData && (
        <div
          className="Search-Tree"
          ref={treeContainerRef}
          style={{ width: '100%', height: 'auto' }}
        >
          <h2 className="Search-subtitle">HASIL PENCARIAN</h2>
          <Tree
            data={treeData}
            orientation="vertical"
            translate={{ x: dimensions.width / 2, y: 50 }}
            nodeSize={{ x: 70, y: 80 }}
            zoomable={true}
            pathFunc="diagonal"
            separation={{ siblings: 1.5, nonSiblings: 2 }}
            renderCustomNodeElement={renderCustomNode}
          />
        </div>
      )}
    </div>
  );
};

export default About;
