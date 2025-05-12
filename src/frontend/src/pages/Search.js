import { React, useState, useRef, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import Tree from 'react-d3-tree';
import BFS from '../media/icons/BFS.svg';
import DFS from '../media/icons/DFS.svg';
import single from '../media/icons/single.svg';
import multiple from '../media/icons/multiple.svg';

const About = () => {
  const { register, handleSubmit } = useForm();
  const [treeDataList, setTreeDataList] = useState([]);
  const [metaInfo, setMetaInfo] = useState({ timetaken: "-", node_visited: 0 });
  const treeContainerRef = useRef(null);
  const [containerWidth, setContainerWidth] = useState(0);
  const [selectedAlgorithm, setSelectedAlgorithm] = useState(null);
  const [selectedSearchMode, setSelectedSearchMode] = useState(null);

  const onSubmit = async (querySearch) => {
    querySearch.maksimalResep = Number(querySearch.maksimalResep);
    console.log("Kirim ke backend:", querySearch); // debug di console (liat inspect)

    try {
      const response = await fetch("http://localhost:8080/api/search", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(querySearch),
      });

      const data = await response.json();

      if (data.treeData && Array.isArray(data.treeData)) {
        setTreeDataList(data.treeData);
        setMetaInfo({
          timetaken: data.timetaken || "-",
          node_visited: data.node_visited || 0,
        });
      } else {
        console.error("Unexpected API response format:", data);
      }
    } catch (err) {
      console.error("Failed to fetch:", err);
    }
  };

  useEffect(() => {
    if (treeContainerRef.current) {
      setContainerWidth(treeContainerRef.current.offsetWidth);
    }
  }, [treeDataList]);

  const calculateTreeDepth = (node) => {
    if (!node.children || node.children.length === 0) return 1;
    return 1 + Math.max(...node.children.map(calculateTreeDepth));
  };

  const calculateTreeWidth = (node, levelWidths = {}, depth = 0) => {
    if (!levelWidths[depth]) levelWidths[depth] = 0;
    levelWidths[depth] += 1;
    if (node.children) {
      node.children.forEach(child => calculateTreeWidth(child, levelWidths, depth + 1));
    }
    return Math.max(...Object.values(levelWidths));
  };
  const renderCustomNode = ({ nodeDatum }) => {
    const isLeaf = !nodeDatum.children || nodeDatum.children.length === 0;
    const textLength = nodeDatum.name.length;
    const width = 80;
    const baseHeight = 40;
    const estimatedLineCount = Math.ceil(textLength / 16);
    const height = baseHeight * estimatedLineCount;
    return (
      <g>
        <foreignObject width={width} height={height} x={-width / 2} y={-height / 2}>
          <div
            xmlns="http://www.w3.org/1999/xhtml"
            className={`tree-node ${isLeaf ? 'leaf-node' : 'parent-node'}`}
          >
            <p className="node-label-multiline">{nodeDatum.name}</p>
          </div>
        </foreignObject>
      </g>
    );
  };

  const OptionsButton = ({ Opsi1, Icon1, Opsi2, Icon2, dataType, onOptionSelect }) => {
    const [selectedOption, setSelectedOption] = useState(null);
    const handleClick = (value) => {
      setSelectedOption(value);
      onOptionSelect(value);
    };

    return (
    <div className="options-container">
      <button
        type="button"
        className={`options-button ${selectedOption === Opsi1 ? 'active' : ''}`}
        onClick={() => handleClick(Opsi1)}
      >
          <img src={Icon1} alt={Opsi1} />
          {Opsi1}
      </button>
      <button
        type="button"
        className={`options-button ${selectedOption === Opsi2 ? 'active' : ''}`}
        onClick={() => handleClick(Opsi2)}
      >
          <img src={Icon2} alt={Opsi2} />
          {Opsi2}
      </button>

      <input
        type="hidden"
        value={selectedOption}
        {...register(dataType, { required: true })}
      />
    </div>
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
              placeholder="Contoh: Brick, Cloud"
              autoComplete="off"
              className="custom-search-input"
              {...register("namaResep", { required: true })}
            />
          </div>

          <div className="Search-form-card">
            <h3>Jumlah Resep*</h3>
            <input
              type="number"
              placeholder="Contoh: 5"
              className="custom-search-input"
              {...register("maksimalResep", { required: true, min: 1 })}
            />
          </div>

          <div className="Search-form-card">
            <h3>Algoritma*</h3>
            <OptionsButton Opsi1 = 'BFS' Opsi2 = 'DFS' Icon1 = {BFS} Icon2 = {DFS} dataType='algoritma'
            onOptionSelect={(value) => setSelectedAlgorithm(value)} />
          </div>

          <div className="Search-form-card">
            <h3>Mode Pencarian*</h3>
            <OptionsButton Opsi1 = 'Single' Opsi2 = 'Multiple' Icon1 = {single} Icon2 = {multiple} dataType='modePencarian'
            onOptionSelect={(value) => setSelectedSearchMode(value)}/>
          </div>
          
          {selectedAlgorithm === 'DFS' && selectedSearchMode === 'Multiple' && (
          <div className="Search-form-card">
            <h3>Mode Multiple DFS*</h3>
            <div className="Search-form-card-radio">
              <div className="Search-form-card-radio-2">
                <label className="custom-radio">
                  <input
                    type="radio"
                    {...register("modeMultipleDFS", { required: true })}
                    value="S" /* Sequence */
                  />
                  <span className="radio-image" />
                </label>

                <label className="custom-radio">
                  <input
                    type="radio"
                    {...register("modeMultipleDFS", { required: true })}
                    value="SM" /* Sequence Mutex */
                  />
                  <span className="radio-image" />
                </label>

                <label className="custom-radio">
                  <input
                    type="radio"
                    {...register("modeMultipleDFS", { required: true })}
                    value="ELP" /* Element Level Parallelism */
                  />
                  <span className="radio-image" />
                </label>

              </div>
              <div className="Search-form-card-radio-3">
                <p>Pure Sequence</p>
                <p>Sequence + Mutex</p>
                <p>Element Level Parallelism</p>
              </div>
            </div>
          </div>
          )}

          {selectedAlgorithm === 'BFS' && selectedSearchMode === 'Multiple' && (
          <div className="Search-form-card">
            <h3>Mode Multiple BFS*</h3>
            <div className="Search-form-card-radio">
              <div className="Search-form-card-radio-2">
                <label className="custom-radio">
                  <input
                    type="radio"
                    {...register("modeMultipleBFS", { required: true })}
                    value="ApaNau" /* Sequence */
                  />
                  <span className="radio-image" />
                </label>

                <label className="custom-radio">
                  <input
                    type="radio"
                    {...register("modeMultipleBFS", { required: true })}
                    value="SM" /* Sequence Mutex */
                  />
                  <span className="radio-image" />
                </label>

                <label className="custom-radio">
                  <input
                    type="radio"
                    {...register("modeMultipleBFS", { required: true })}
                    value="ELP" /* Element Level Parallelism */
                  />
                  <span className="radio-image" />
                </label>

              </div>
              <div className="Search-form-card-radio-3">
                <p>Apa Nau</p>
                <p>Sequence + Mutex</p>
                <p>Element Level Parallelism</p>
              </div>
            </div>
          </div>
          )}

          

          <input type="submit" className="submit-button" />
        </div>
      </form>

      {/* Tree View */}
      {treeDataList.length > 0 && (
        <div className="Search-Tree" ref={treeContainerRef} style={{ width: '100%', minHeight: '100vh' }}>
          <h1 className="Search-title">HASIL PENCARIAN</h1>
          <div className="Search-meta-info">
            <p>Waktu Pencarian: {metaInfo.timetaken}</p>
            <p>Node yang Dikunjungi: {metaInfo.node_visited}</p>
          </div>
          {treeDataList.map((treeData, idx) => {
            const depth = calculateTreeDepth(treeData);
            const height = Math.max(depth * 120, 300);
            const treeWidth = calculateTreeWidth(treeData) * 120;
            return (
              <div key={idx} style={{ height: `${height}px`}}>
                <h3 className="Search-subtitle-2">{treeData.name}</h3>
                <Tree
                  data={treeData}
                  orientation="vertical"
                  translate={{ x: containerWidth / 2, y: 50 }}
                  nodeSize={{ x: 90, y: 70 }}
                  zoomable={false}
                  initialZoom={Math.min(containerWidth / treeWidth, 1)}
                  pathFunc="diagonal"
                  separation={{ siblings: 1.2, nonSiblings: 1.5 }}
                  renderCustomNodeElement={renderCustomNode}
                />
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default About;
