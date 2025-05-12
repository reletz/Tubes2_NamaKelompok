import React from 'react';
import { useForm } from 'react-hook-form';
import Tree from 'react-d3-tree';
import { parseMultipleTrees, parseMetaInfo } from './parseTree';
import rawTree from '../data/multi_dfs_results.json';

const About = () => {
  const { register, handleSubmit } = useForm();
  const [treeDataList, setTreeDataList] = React.useState([]);
  const [metaInfo, setMetaInfo] = React.useState({ timetaken: "-", node_visited: 0 });
  const treeContainerRef = React.useRef(null);
  const [containerWidth, setContainerWidth] = React.useState(0);

  const onSubmit = (querySearch) => {
    console.log(querySearch);
    const trees = parseMultipleTrees(rawTree);
    const info = parseMetaInfo(rawTree);
    setTreeDataList(trees);
    setMetaInfo(info);
  };

  React.useEffect(() => {
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
              autoComplete="off"
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
              {...register("Maksimal Resep", { required: true, min: 1 })}
            />
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

          <input type="submit" className="submit-button" />
        </div>
      </form>

      {/* Tree View */}
      {treeDataList.length > 0 && (
        <div className="Search-Tree" ref={treeContainerRef} style={{ width: '100%', minHeight: '100vh' }}>
          <h2 className="Search-subtitle">HASIL PENCARIAN</h2>
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
