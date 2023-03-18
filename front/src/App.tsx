import React, { useEffect, useState } from 'react';
import './App.css';
import { ExecuteMtrRequest, ExecuteMtrResponse, ReportHub } from './gen/eyeroute/mtr/v1alpha1/mtr_pb';
import { MtrService } from './gen/eyeroute/mtr/v1alpha1/mtr_connect';
import { createCallbackClient, createPromiseClient } from '@bufbuild/connect';
import { createConnectTransport } from '@bufbuild/connect-web';

function App() {
  const [ipAddress, setIpAddress] = useState("")
  const [reportHubs, setReportHubs] = useState<ReportHub[]>([])

  const client = createCallbackClient(
    MtrService,
    createConnectTransport({
      baseUrl: "http://127.0.0.1:8080/api",
    })
  )

  const onClickExecuteButton = () => {
    if (ipAddress === "") {
      console.error("IP Address is empty!")
    }

    const executeMtrReq = new ExecuteMtrRequest({
      ipAddress: ipAddress,
      reportCycles: 10,
    })

    client.executeMtr(executeMtrReq, (err, res) => {
      if (err) {
        console.error(err)
        return
      }

      console.log(res)
      setReportHubs(res.hubs)
    })
  }

  return (
    <div className="App">
      <h1>eyeroute (Looking Glass)</h1>

      <h2>mtr</h2>
      <input onChange={(e) => setIpAddress(e.target.value)} />
      <button onClick={onClickExecuteButton}>実行</button>

      <table>
        <thead>
          <tr>
            <th>count</th>
            <th>host</th>
            <th>loss</th>
            <th>sent</th>
            <th>last</th>
            <th>avg</th>
            <th>best</th>
            <th>worst</th>
            <th>stdev</th>
          </tr>
        </thead>
        <tbody>
          {reportHubs.map((hub, index) => (
            <tr key={index}>
              <td>{hub.count}</td>
              <td>{hub.host}</td>
              <td>{hub.loss}</td>
              <td>{hub.sent}</td>
              <td>{hub.last}</td>
              <td>{hub.avg}</td>
              <td>{hub.best}</td>
              <td>{hub.worst}</td>
              <td>{hub.stdev}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default App;
