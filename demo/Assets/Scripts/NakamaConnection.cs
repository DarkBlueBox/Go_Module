using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Nakama;
using System.Threading.Tasks;
using System;

[Serializable]
[CreateAssetMenu]
public class NakamaConnection : ScriptableObject
{
    public string Scheme = "http";
    public string Host = "localhost";
    public int Port = 7350;
    public string ServerKey = "defaultkey";

    private const string SessionPrefName = "nakama.session";
    private const string DeviceIdentifierPrefName = "nakama.deviceUniqueIdentifier";

    public IClient Client;
    public ISession Session;
    public ISocket Socket;

    public ISocket GetSocket() { return Socket; }

    private IMatch currentMatch;

    private string currentMatchmakingTicket;
    private string currentMatchId;

    /// <summary>
    /// Connects to the Nakama server using device authentication and opens socket for realtime communication.
    /// </summary>
    public async Task Connect()
    {
        // Connect to the Nakama server.
        Client = new Nakama.Client(Scheme, Host, Port, ServerKey, UnityWebRequestAdapter.Instance);

        // Attempt to restore an existing user session.
        var authToken = PlayerPrefs.GetString(SessionPrefName);
        if (!string.IsNullOrEmpty(authToken))
        {
            var session = Nakama.Session.Restore(authToken);
            if (!session.IsExpired)
            {
                Session = session;
            }
        }

        // If we weren't able to restore an existing session, authenticate to create a new user session.
        if (Session == null)
        {
            string deviceId;

            // If we've already stored a device identifier in PlayerPrefs then use that.
            if (PlayerPrefs.HasKey(DeviceIdentifierPrefName))
            {
                deviceId = PlayerPrefs.GetString(DeviceIdentifierPrefName);
            }
            else
            {
                // If we've reach this point, get the device's unique identifier or generate a unique one.
                deviceId = SystemInfo.deviceUniqueIdentifier;
                if (deviceId == SystemInfo.unsupportedIdentifier)
                {
                    deviceId = System.Guid.NewGuid().ToString();
                }

                // Store the device identifier to ensure we use the same one each time from now on.
                PlayerPrefs.SetString(DeviceIdentifierPrefName, deviceId);
            }

            // Use Nakama Device authentication to create a new session using the device identifier.
            Session = await Client.AuthenticateDeviceAsync(deviceId);

            // Store the auth token that comes back so that we can restore the session later if necessary.
            PlayerPrefs.SetString(SessionPrefName, Session.AuthToken);
        }

        // Open a new Socket for realtime communication.
        Socket = Client.NewSocket();
        await Socket.ConnectAsync(Session, true);

        Debug.Log(Session);
        Debug.Log(Socket);

        
    }
    
    public async Task FindMatch(int minPlayers, int maxPlayers){


        var matchMakingTicket = await Socket.AddMatchmakerAsync("*", minPlayers, maxPlayers);
        currentMatchmakingTicket = matchMakingTicket.Ticket;

    }


    public async Task CancelMatchmaking()
    {
        Debug.Log("Cancel Matchmaking");

        await Socket.RemoveMatchmakerAsync(currentMatchmakingTicket);
    }
    public async Task LeaveMatch()
    {
        Debug.Log("Leave Match");

        await Socket.LeaveMatchAsync(currentMatch.Id);
    }

    public void SetCurrentMatch(IMatch currentMatch)
    {
        this.currentMatch = currentMatch;
    }
}
