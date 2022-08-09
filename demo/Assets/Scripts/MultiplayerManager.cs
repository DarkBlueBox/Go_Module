
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using System.Linq;
using System.Threading.Tasks;
using Nakama;
using Nakama.TinyJson;
public class MultiplayerManager : MonoBehaviour
{

    int minPlayers = 2;
    int maxPlayers = 2;

    public NakamaConnection NakamaConnection;
    public GameObject NetworkLocalPlayerPrefab;
    public GameObject NetworkRemotePlayerPrefab;
    public GameObject MainMenu;
    public GameObject SpawnPoints;


    private IDictionary<string, GameObject> players;
    private IUserPresence localUser;
    private GameObject localPlayer;
    private IMatch currentMatch;

    private Transform[] spawnPoints;

    private string localDisplayName;

    [SerializeField] public string localUserSessionID;

    private async void Start()
    {

        players = new Dictionary<string, GameObject>();


        var mainThread = UnityMainThreadDispatcher.Instance();

        await NakamaConnection.Connect();


        NakamaConnection.Socket.ReceivedMatchmakerMatched += m => mainThread.Enqueue(() => OnReceivedMatchmakerMatched(m));
        NakamaConnection.Socket.ReceivedMatchPresence += m => mainThread.Enqueue(() => OnReceivedMatchPresence(m));

      
    }

      public async Task FindMatch()
    {
        await NakamaConnection.FindMatch(minPlayers, maxPlayers);
    }

    public async Task CanelMatchMacking()
    {
        await NakamaConnection.CancelMatchmaking();
    }

    public async Task LeaveMatch()
    {
        
        string jsonState = MatchData.SetUserID(localUserSessionID);
        await SendMatchStateAsync(OpCodes.PlayerLeft, jsonState);

        await NakamaConnection.LeaveMatch();
    }


     public async Task HelloWorld()
    {
        
        await NakamaConnection.Socket.SendMatchStateAsync(currentMatch.Id, OpCodes.Hello,"",null);
       
    }

    private async void OnReceivedMatchmakerMatched(IMatchmakerMatched matched)
    {

        localUser = matched.Self.Presence;

        var match = await NakamaConnection.Socket.JoinMatchAsync(matched);

        
        Debug.Log("Session id" + match.Self.SessionId);

   

        foreach (var user in match.Presences)
        {
          Debug.Log("user Session id" + user.SessionId);
        }


        currentMatch = match;
    }

    private void OnReceivedMatchPresence(IMatchPresenceEvent matchPresenceEvent)
    {


    }


    private async Task OnReceivedMatchState(IMatchState matchState)
    {

        var userSessionId = matchState.UserPresence.SessionId;


        var state = matchState.State.Length > 0 ? System.Text.Encoding.UTF8.GetString(matchState.State).FromJson<Dictionary<string, string>>() : null;

        switch(matchState.OpCode)
        {
           case OpCodes.Hello:
                Debug.Log("Hello");
                await NakamaConnection.Socket.SendMatchStateAsync(currentMatch.Id,OpCodes.World,"",new []{matchState.UserPresence});
                break;
            case OpCodes.World:
                Debug.Log("World");
                await NakamaConnection.Socket.SendMatchStateAsync(currentMatch.Id,OpCodes.Hello,"",new []{matchState.UserPresence});
                break;
        }
    }



    public async Task SendMatchStateAsync(long opCode, string state)
    {
        await NakamaConnection.Socket.SendMatchStateAsync(currentMatch.Id, opCode, state);
    }

    public void SendMatchState(long opCode, string state)
    {
        NakamaConnection.Socket.SendMatchStateAsync(currentMatch.Id, opCode, state);
    }
}