/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * ER Diagram Editor Page
 */

import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { ERDiagramList } from '@coze-extensions/er-diagram';
import type { ERDiagram } from '@coze-extensions/er-diagram';

export function ERDiagramPage() {
  const { space_id } = useParams<{ space_id: string }>();
  const workspaceId = space_id ? parseInt(space_id, 10) : 0;
  const [selectedDiagram, setSelectedDiagram] = useState<ERDiagram | null>(null);

  if (!workspaceId) {
    return <div>无效的工作空间 ID</div>;
  }

  const handleDiagramSelect = (diagram: ERDiagram) => {
    setSelectedDiagram(diagram);
    // TODO: 打开 ER 图编辑器
    console.log('Selected diagram:', diagram);
  };

  const handleCreateDiagram = (diagram: ERDiagram) => {
    setSelectedDiagram(diagram);
    // TODO: 打开 ER 图编辑器
    console.log('Created diagram:', diagram);
  };

  return (
    <div className="er-diagram-page">
      <ERDiagramList
        workspaceId={workspaceId}
        onDiagramSelect={handleDiagramSelect}
        onCreateDiagram={handleCreateDiagram}
      />
      {selectedDiagram && (
        <div className="er-diagram-editor-placeholder">
          <p>ER 图编辑器（待实现可视化编辑器）</p>
          <p>Diagram ID: {selectedDiagram.id}</p>
          <p>Name: {selectedDiagram.name}</p>
        </div>
      )}
    </div>
  );
}

