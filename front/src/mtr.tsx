import { createCallbackClient } from '@bufbuild/connect';
import { createConnectTransport } from '@bufbuild/connect-web';
import { useState } from 'react';
import { MtrService } from './gen/eyeroute/mtr/v1alpha1/mtr_connect';
import { ExecuteMtrRequest, ReportHub } from './gen/eyeroute/mtr/v1alpha1/mtr_pb';

function Mtr() {
  const [ipAddress, setIpAddress] = useState("")
  const [reportHubs, setReportHubs] = useState<ReportHub[]>([])

  const client = createCallbackClient(
    MtrService,
    createConnectTransport({
      baseUrl: "/",
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
    <div>
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

export default Mtr;
