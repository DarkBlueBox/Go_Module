using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Nakama.TinyJson;

public class MatchData : MonoBehaviour
{
     public static string SetPlayerState(string userSessionID, string playerState)
    {
        var values = new Dictionary<string, string>
        {
            { "userSessionID", userSessionID },
            { "playerState", playerState },
        };

        return values.ToJson();
    }
    public static string Died(Vector3 position)
    {
        var values = new Dictionary<string, string>
        {
            { "position.x", position.x.ToString() },
            { "position.y", position.y.ToString() }
        };

        return values.ToJson();
    }


    public static string Respawned(int spawnIndex)
    {
        var values = new Dictionary<string, string>
        {
            { "spawnIndex", spawnIndex.ToString() },
        };

        return values.ToJson();
    }

    public static string SetUserID(string userSessionID)
    {
        var values = new Dictionary<string, string>
        {
            { "userSessionID", userSessionID }
        };

        return values.ToJson();
    }
}
