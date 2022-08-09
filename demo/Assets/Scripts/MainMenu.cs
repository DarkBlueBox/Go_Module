using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class MainMenu : MonoBehaviour
{


    MultiplayerManager multiplayermanager;


    void Start()
    {
        multiplayermanager = FindObjectOfType<MultiplayerManager>();

    }

    public async void FindMatch()
    {
        await multiplayermanager.FindMatch();
    }

    public async void LeaveMatch()
    {
        await multiplayermanager.LeaveMatch();
    }
    public async void CancelMatchMaking()
    {
        await multiplayermanager.CanelMatchMacking();
    }

     public async void HelloWorld()
    {
        await multiplayermanager.HelloWorld();
    }
    
}
